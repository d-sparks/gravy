package dailyprices

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/bradfitz/slice"
	"github.com/d-sparks/gravy/data/alpha"
	"github.com/d-sparks/gravy/data/covariance"
	dailyprices_pb "github.com/d-sparks/gravy/data/dailyprices/proto"
	"github.com/d-sparks/gravy/data/mean"
	"github.com/d-sparks/gravy/data/variance"
	"github.com/d-sparks/gravy/gravy"
	"github.com/golang/protobuf/ptypes"
	_ "github.com/lib/pq"
)

// Stats are the values tracked per ticker.
type Stats struct {
	alpha                *alpha.Rolling
	movingMean           *mean.Rolling
	movingMeanReturn     *mean.Rolling
	movingVariance       map[int]*variance.Rolling
	movingMeanVolume     *mean.Rolling
	movingVolumeVariance map[int]*variance.Rolling
	mean                 *mean.Streaming
	variance             *variance.Streaming
}

// PairStats are things we measure for pairs of tickers.
type PairStats struct {
	covariance *covariance.Streaming
}

// Server implements dailyprices_pb.DataServer.
type Server struct {
	dailyprices_pb.UnimplementedDataServer

	// PostGRES
	db                    *sql.DB
	pricesTableName       string
	tradingDatesTableName string
	assetIDsTableName     string

	// Cache
	mu        sync.Mutex
	cache     map[int32]map[time.Time]*dailyprices_pb.DailyData
	times     []time.Time
	timeIndex map[time.Time]int

	// Measurements
	stats     map[string]*Stats
	pairStats map[string]map[string]*PairStats

	// Tracking
	firstSeen    map[string]time.Time
	lastSeen     map[string]time.Time
	missingDates map[string]map[time.Time]struct{}

	// Benchmark for the market. (Currently SPY)
	benchmark *mean.Rolling

	// Bad dates (hack): Skip trading on these days.
	badDates map[time.Time]struct{}
}

// NewServer creates an empty daily prices server.
func NewServer(
	postgresURL string,
	dailyPricesTable string,
	tradingDatesTable string,
	assetIDsTable string,
) (*Server, error) {
	log.Printf("Connecting to database `%s`", postgresURL)
	db, err := sql.Open("postgres", postgresURL)
	if err != nil {
		return nil, fmt.Errorf("Error connecting to postgres: %s", err.Error())
	}

	unsafeParse := func(date string) time.Time {
		timeStr, _ := time.Parse("2006-01-02", date)
		return timeStr
	}
	server := Server{
		db:                    db,
		pricesTableName:       dailyPricesTable,
		tradingDatesTableName: tradingDatesTable,
		assetIDsTableName:     assetIDsTable,
		cache:                 map[int32]map[time.Time]*dailyprices_pb.DailyData{},
		badDates: map[time.Time]struct{}{
			unsafeParse("2011-02-17"): struct{}{},
		},
	}

	go server.RunDebugServer(8080)

	return &server, nil
}

// Close the underlying postgres connection.
func (s *Server) Close() {
	s.db.Close()
}

// updateStatsForTicker updates the averages and alpha and returns the relative performance.
func (s *Server) updateStatsForTicker(
	tickTime time.Time,
	ticker string,
	benchmarkPerf float64,
	data *dailyprices_pb.DailyData,
) float64 {
	stats := s.stats[ticker]

	// Get the new price.
	prices := data.GetPrices()[ticker]
	price := prices.GetClose()
	volume := prices.GetVolume()

	// Get relative performance. If ticker is the benchmark asset, set benchmarkPerf.
	perf := gravy.RelativePerfOrZero(price, stats.movingMean.OldestObservation())
	if ticker == "SPY" {
		benchmarkPerf = perf
	}

	// Update moving means first.
	stats.movingMean.Observe(price)
	stats.movingMeanReturn.Observe(perf)
	stats.movingMeanVolume.Observe(volume)

	// Now update moving variances.
	stats.alpha.Observe(perf, benchmarkPerf)
	stats.movingVariance[15].Observe(price)
	stats.movingVariance[35].Observe(price)
	stats.movingVariance[252].Observe(price)
	stats.movingVolumeVariance[15].Observe(volume)

	// Update streaming quantities.
	stats.mean.Observe(price)
	stats.variance.Observe(price)

	// Update dailyPrices proto.
	if _, ok := data.GetStats()[ticker]; !ok {
		data.Stats[ticker] = &dailyprices_pb.Stats{
			Exchange: s.cache[0][s.firstSeen[ticker]].GetStats()[ticker].GetExchange(),
		}
	}
	out := data.Stats[ticker]
	out.Alpha = stats.alpha.Alpha()
	out.Beta = stats.alpha.Beta()
	out.MovingAverages = map[int32]float64{
		15:  stats.movingMean.Value(15),
		35:  stats.movingMean.Value(35),
		252: stats.movingMean.Value(252),
	}
	out.MovingAverageReturns = map[int32]float64{
		15:  stats.movingMeanReturn.Value(15),
		35:  stats.movingMeanReturn.Value(35),
		252: stats.movingMeanReturn.Value(252),
	}
	out.MovingVariance = map[int32]float64{
		15:  stats.movingVariance[15].Value(),
		35:  stats.movingVariance[35].Value(),
		252: stats.movingVariance[252].Value(),
	}
	out.MovingVolume = map[int32]float64{
		15: stats.movingMeanVolume.Value(15),
	}
	out.MovingVolumeVariance = map[int32]float64{
		15: stats.movingVolumeVariance[15].Value(),
	}
	out.Mean = stats.mean.Value()
	out.Variance = stats.variance.Value()

	return perf
}

// validTickerPair returns true if the first string is less than the second string and the second string is not equal to
// the first string + "_".
func validTickerPair(first, second string) bool {
	return first < second && second != first+"_"
}

// updateStatsForTickerPair updates the stats for the given ticker pair.
func (s *Server) updateStatsForTickerPair(
	tickTime time.Time,
	ticker string,
	otherTicker string,
	data *dailyprices_pb.DailyData,
) {
	tickerPrice := data.GetPrices()[ticker].GetClose()
	otherTickerPrice := data.GetPrices()[otherTicker].GetClose()
	if _, ok := s.pairStats[ticker][otherTicker]; !ok {
		s.pairStats[ticker][otherTicker] = &PairStats{covariance: covariance.NewStreaming()}
	}
	s.pairStats[ticker][otherTicker].covariance.Observe(tickerPrice, otherTickerPrice)
}

// buildOutputForTickerPairs builds output for ticker pairs. In particular, calculates the correlation and sorts the
// output by descending correlation.
func (s *Server) buildOutputForTickerPairs(data *dailyprices_pb.DailyData) {
	allPairs := []*dailyprices_pb.PairStats{}
	for ticker := range s.firstSeen {
		for otherTicker := range s.firstSeen {
			if !validTickerPair(ticker, otherTicker) {
				continue
			}
			covariance := s.pairStats[ticker][otherTicker].covariance
			if covariance.NumObservations() < 100.0 {
				continue
			}
			allPairs = append(allPairs, &dailyprices_pb.PairStats{
				First:       ticker,
				Second:      otherTicker,
				Covariance:  covariance.Value(),
				Correlation: covariance.CorrelationValue(),
			})
		}
	}

	slice.Sort(allPairs, func(i, j int) bool { return allPairs[i].Correlation > allPairs[j].Correlation })

	if len(allPairs) >= 2000 {
		allPairs = append(allPairs[:1000], allPairs[len(allPairs)-1000:]...)
	}
	data.PairStats = make([]*dailyprices_pb.PairStats, len(allPairs))
	copy(data.PairStats, allPairs)
}

// updateStats updates the various tracked rolling averages at the given time. Also updates the given prices pointer.
func (s *Server) updateStats(tickTime time.Time, data *dailyprices_pb.DailyData) error {
	// Update firstSeen if this is the first time seeing this ticker.
	for ticker := range data.GetPrices() {
		if _, ok := s.firstSeen[ticker]; !ok {
			s.stats[ticker] = &Stats{
				movingMean:       mean.NewRolling(15, 35, 252),
				movingMeanReturn: mean.NewRolling(15, 35, 252),
				movingMeanVolume: mean.NewRolling(15),
				mean:             mean.NewStreaming(),
			}
		}
	}

	// If this is the first day, find and track the benchmark (SPY).
	if tickTime == s.times[0] {
		if stats, ok := s.stats["SPY"]; !ok {
			return fmt.Errorf("Error: Could not find market benchmark.")
		} else {
			s.benchmark = stats.movingMeanReturn
		}
	}

	// Now create alphas if necessary.
	for ticker := range data.GetPrices() {
		if _, ok := s.firstSeen[ticker]; !ok {
			s.firstSeen[ticker] = tickTime
			stats := s.stats[ticker]
			stats.alpha = alpha.NewRolling(stats.movingMeanReturn, s.benchmark, 252, 0.0)
			stats.movingVariance = map[int]*variance.Rolling{
				15:  variance.NewRolling(stats.movingMean, 15),
				35:  variance.NewRolling(stats.movingMean, 35),
				252: variance.NewRolling(stats.movingMean, 252),
			}
			stats.movingVolumeVariance = map[int]*variance.Rolling{
				15: variance.NewRolling(stats.movingMeanVolume, 15),
			}
			stats.variance = variance.NewStreaming()
			s.pairStats[ticker] = map[string]*PairStats{}
		}
	}

	// Update averages first for SPY, the benchmark, as it is necessary for the comparisons that are made to it.
	benchmarkPerf := s.updateStatsForTicker(tickTime, "SPY", 0.0, data)
	for ticker := range s.firstSeen {
		if ticker == "SPY" {
			continue
		}
		s.updateStatsForTicker(tickTime, ticker, benchmarkPerf, data)
	}
	for ticker := range s.firstSeen {
		for otherTicker := range s.firstSeen {
			if !validTickerPair(ticker, otherTicker) {
				continue
			}
			s.updateStatsForTickerPair(tickTime, ticker, otherTicker, data)
		}
	}

	// Builds the output for ticker pairs.
	s.buildOutputForTickerPairs(data)

	return nil
}

// Get implements the get endpoint for dailyprices_pb.DataServer.
func (s *Server) Get(ctx context.Context, req *dailyprices_pb.Request) (*dailyprices_pb.DailyData, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Parse timestamp to Golang native time.
	tickTime, err := ptypes.Timestamp(req.GetTimestamp())
	if err != nil {
		return nil, fmt.Errorf("Invalid timestamp: %s", err.Error())
	}

	// Check cache.
	if versionPrices, ok := s.cache[req.GetVersion()]; ok {
		if cachedDailyData, ok := versionPrices[tickTime]; ok {
			return cachedDailyData, nil
		}
	}

	// Query database.
	rows, err := s.db.Query(
		fmt.Sprintf(
			"SELECT ticker, open, close, low, high, volume, exchange FROM %s WHERE date = $1",
			s.pricesTableName,
		),
		tickTime.Format("2006-01-02"),
	)
	if err != nil {
		return nil, fmt.Errorf("Error reading from db: `%s`", err.Error())
	}
	defer rows.Close()

	// Construct daily prices by scanning the query result.
	data := dailyprices_pb.DailyData{
		Prices: map[string]*dailyprices_pb.Prices{},
		Stats:  map[string]*dailyprices_pb.Stats{},
	}
	for rows.Next() {
		var prices dailyprices_pb.Prices
		var ticker string
		var exchange string
		err := rows.Scan(
			&ticker,
			&prices.Open,
			&prices.Close,
			&prices.Low,
			&prices.High,
			&prices.Volume,
			&exchange,
		)
		if err != nil {
			return nil, fmt.Errorf("Error while parsing row: `%s`", err.Error())
		}
		data.Prices[ticker] = &prices
		data.Stats[ticker] = &dailyprices_pb.Stats{Exchange: exchange}
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("Error constructing response: `%s`", rows.Err().Error())
	}

	// Update averages.
	if err = s.updateStats(tickTime, &data); err != nil {
		return nil, fmt.Errorf("Error updating averages: %s", err.Error())
	}

	// Stamp, cache, and return.
	data.Timestamp = req.GetTimestamp()
	data.Version = req.GetVersion()
	if _, ok := s.cache[data.GetVersion()]; !ok {
		s.cache[data.GetVersion()] = map[time.Time]*dailyprices_pb.DailyData{}
	}
	s.cache[data.GetVersion()][tickTime] = &data
	return &data, nil
}

// NewSession implements the interface and sets/resets state for the session.
func (s *Server) NewSession(
	ctx context.Context,
	req *dailyprices_pb.NewSessionRequest,
) (*dailyprices_pb.NewSessionResponse, error) {
	// Get trading dates for entire session.
	dates, err := s.TradingDatesInRange(ctx, req.GetSimRange())
	if err != nil {
		return nil, fmt.Errorf("Error getting dates to start session: %s", err.Error())
	}
	s.times = make([]time.Time, len(dates.GetTimestamps()))
	s.timeIndex = map[time.Time]int{}
	for i, date := range dates.GetTimestamps() {
		if s.times[i], err = ptypes.Timestamp(date); err != nil {
			return nil,
				fmt.Errorf("Invalid date/timestamp: %s", err.Error())
		}
		s.timeIndex[s.times[i]] = i
	}

	// Reset measurements.
	s.stats = map[string]*Stats{}
	s.firstSeen = map[string]time.Time{}
	s.lastSeen = map[string]time.Time{}
	s.missingDates = map[string]map[time.Time]struct{}{}
	s.pairStats = map[string]map[string]*PairStats{}

	return &dailyprices_pb.NewSessionResponse{}, nil
}

// TradingDatesInRange implements the interface method. Returns trading dates in the given range.
func (s *Server) TradingDatesInRange(
	ctx context.Context,
	dateRange *dailyprices_pb.Range,
) (*dailyprices_pb.TradingDates, error) {
	// Parse timestamps to Golang native time.
	lb, err := ptypes.Timestamp(dateRange.GetLb())
	if err != nil {
		return nil, fmt.Errorf("Invalid lb timestamp: %s", err.Error())
	}
	ub, err := ptypes.Timestamp(dateRange.GetUb())
	if err != nil {
		return nil, fmt.Errorf("Invalid ub timestamp: %s", err.Error())
	}

	// Query for dates.
	rows, err := s.db.Query(
		fmt.Sprintf(
			"SELECT DISTINCT date FROM %s WHERE date >= $1 AND date <= $2 ORDER BY date",
			s.tradingDatesTableName,
		),
		lb.Format("2006-01-02"),
		ub.Format("2006-01-02"),
	)
	if err != nil {
		return nil, fmt.Errorf("Error querying for distinct dates: `%s`", err.Error())
	}
	defer rows.Close()

	// Scan and parse dates into a slice.
	tradingDates := dailyprices_pb.TradingDates{}
	for rows.Next() {
		var dateStr string
		if err := rows.Scan(&dateStr); err != nil {
			return nil, fmt.Errorf("Error scanning date `%s` from postgres: %s", dateStr, err.Error())
		}
		date, err := time.Parse("2006-01-02T15:04:05Z", dateStr)
		if err != nil {
			return nil, fmt.Errorf("Could not parse date `%s`: %s", dateStr, err.Error())
		}
		// Skip bad dates: This is a hack to account for questionable data.
		if _, ok := s.badDates[date]; ok {
			continue
		}
		dateProto, err := ptypes.TimestampProto(date)
		if err != nil {
			return nil, fmt.Errorf("Error parsing timestamp to proto: %s", err.Error())
		}
		tradingDates.Timestamps = append(tradingDates.Timestamps, dateProto)

	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("Error scanning rows for distinct dates: %s", rows.Err().Error())
	}

	return &tradingDates, nil
}

// AssetIds implements the interface. Queries the assetIDs table and returns a map of id to (exchange, ticker) pairs.
func (s *Server) AssetIds(
	ctx context.Context,
	unusedInput *dailyprices_pb.AssetIdsRequest,
) (*dailyprices_pb.AssetIdsResponse, error) {
	// Query
	rows, err := s.db.Query(fmt.Sprintf("SELECT exchange, ticker, id FROM %s;", s.assetIDsTableName))
	if err != nil {
		return nil, fmt.Errorf("Error querying for distinct dates: `%s`", err.Error())
	}

	// Scan
	output := dailyprices_pb.AssetIdsResponse{}
	for rows.Next() {
		var (
			ticker   string
			exchange string
			id       int64
		)
		if err = rows.Scan(&exchange, &ticker, &id); err != nil {
			return nil, fmt.Errorf("Error scanning: %s", err.Error())
		}
		output.AssetIds[id] = &dailyprices_pb.AssetIdsResponse_TickerExchangePair{
			Ticker:   ticker,
			Exchange: exchange,
		}
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("Couldn't complete scan of asset ids: %s", err.Error())
	}
	rows.Close()

	// Return
	return &output, nil
}
