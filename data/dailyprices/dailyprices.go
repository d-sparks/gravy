package dailyprices

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/d-sparks/gravy/data/alpha"
	dailyprices_pb "github.com/d-sparks/gravy/data/dailyprices/proto"
	"github.com/d-sparks/gravy/data/mean"
	"github.com/golang/protobuf/ptypes"
	_ "github.com/lib/pq"
)

// Server implements dailyprices_pb.DataServer.
type Server struct {
	dailyprices_pb.UnimplementedDataServer

	// PostGRES
	db                    *sql.DB
	pricesTableName       string
	tradingDatesTableName string

	// Cache
	mu        sync.Mutex
	cache     map[int32]map[time.Time]*dailyprices_pb.DailyPrices
	times     []time.Time
	timeIndex map[time.Time]int

	// Track 15, 35, and 252 day rolling averages. (~20, 50, 365 days worth of trading days.)
	rollingAverages       map[string]*mean.Rolling
	rollingAverageReturns map[string]*mean.Rolling

	// Alpha for each ticker.
	alpha map[string]*alpha.Rolling

	// Benchmark for the market. (Currently SPY)
	benchmark *mean.Rolling

	// First seen time index for each asset.
	firstSeen    map[string]time.Time
	lastSeen     map[string]time.Time
	missingDates map[string]map[time.Time]struct{}
}

// NewServer creates an empty daily prices server.
func NewServer(
	postgresURL string,
	dailyPricesTable string,
	tradingDatesTable string,
) (*Server, error) {
	log.Printf("Connecting to database `%s`", postgresURL)
	db, err := sql.Open("postgres", postgresURL)
	if err != nil {
		return nil, fmt.Errorf("Error connecting to postgres: %s", err.Error())
	}

	var server Server
	server.db = db
	server.pricesTableName = dailyPricesTable
	server.tradingDatesTableName = tradingDatesTable
	server.cache = map[int32]map[time.Time]*dailyprices_pb.DailyPrices{}

	return &server, nil
}

// Close the underlying postgres connection.
func (s *Server) Close() {
	s.db.Close()
}

// updateAveragesForTicker updates the averages and alpha and returns the relative performance.
func (s *Server) updateAveragesForTicker(
	tickTime time.Time,
	ticker string,
	benchmarkPerf float64,
	dailyPrices *dailyprices_pb.DailyPrices,
) float64 {
	// Get the new price.
	price := 0.0
	if p, ok := dailyPrices.GetStockPrices()[ticker]; ok {
		price = p.GetClose()
	}

	// Get relative performance.
	perf := 0.0
	if price0 := s.rollingAverages[ticker].OldestValue(); price0 > 0.0 {
		perf = (price - price0) / price0
	}

	// Update rolling average.
	s.rollingAverages[ticker].Observe(price)
	s.rollingAverageReturns[ticker].Observe(perf)

	// If this is the benchmark asset, it is the benchmark. (And the given benchmark ought to be 0.0.)
	if ticker == "SPY" {
		benchmarkPerf = perf
	}

	// Assumes the benchmark rollingAverageReturns has been updated already.
	s.alpha[ticker].Observe(perf, benchmarkPerf)

	// Update dailyPrices proto.
	if _, ok := dailyPrices.Measurements[ticker]; !ok {
		dailyPrices.Measurements[ticker] = &dailyprices_pb.Measurements{
			Exchange: s.cache[0][s.firstSeen[ticker]].GetMeasurements()[ticker].GetExchange(),
		}
	}
	dailyPrices.Measurements[ticker].MovingAverages = map[int32]float64{
		15:  s.rollingAverages[ticker].Value(15),
		35:  s.rollingAverages[ticker].Value(35),
		252: s.rollingAverages[ticker].Value(252),
	}
	dailyPrices.Measurements[ticker].Alpha = s.alpha[ticker].Alpha()
	dailyPrices.Measurements[ticker].Beta = s.alpha[ticker].Beta()

	return perf
}

// updateAverages updates the various tracked rolling averages at the given time. Also updates the given prices pointer.
func (s *Server) updateAverages(tickTime time.Time, dailyPrices *dailyprices_pb.DailyPrices) error {
	// Update firstSeen if this is the first time seeing this ticker.
	for ticker := range dailyPrices.GetStockPrices() {
		if _, ok := s.firstSeen[ticker]; !ok {
			s.rollingAverages[ticker] = mean.NewRolling(15, 35, 252)
			s.rollingAverageReturns[ticker] = mean.NewRolling(15, 35, 252)
		}
	}

	// If this is the first day, find and track the benchmark (SPY).
	if tickTime == s.times[0] {
		ok := true
		if s.benchmark, ok = s.rollingAverageReturns["SPY"]; !ok {
			return fmt.Errorf("Error: Could not find market benchmark.")
		}
	}

	// Now create alphas if necessary.
	for ticker := range dailyPrices.GetStockPrices() {
		if _, ok := s.firstSeen[ticker]; !ok {
			s.firstSeen[ticker] = tickTime
			s.alpha[ticker] = alpha.NewRolling(s.rollingAverageReturns[ticker], s.benchmark, 252, 0.0)
		}
	}

	// Update averages first for SPY, the benchmark, as it is necessary for the comparisons that are made to it.
	benchmarkPerf := s.updateAveragesForTicker(tickTime, "SPY", 0.0, dailyPrices)
	for ticker := range s.firstSeen {
		s.updateAveragesForTicker(tickTime, ticker, benchmarkPerf, dailyPrices)
	}

	return nil
}

// Get implements the get endpoint for dailyprices_pb.DataServer.
func (s *Server) Get(ctx context.Context, req *dailyprices_pb.Request) (*dailyprices_pb.DailyPrices, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Parse timestamp to Golang native time.
	tickTime, err := ptypes.Timestamp(req.GetTimestamp())
	if err != nil {
		return nil, fmt.Errorf("Invalid timestamp: %s", err.Error())
	}

	// Check cache.
	if versionPrices, ok := s.cache[req.GetVersion()]; ok {
		if cachedDailyPrices, ok := versionPrices[tickTime]; ok {
			return cachedDailyPrices, nil
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

	// Construct daily prices by scanning the query result.
	var dailyPrices dailyprices_pb.DailyPrices
	dailyPrices.StockPrices = map[string]*dailyprices_pb.DailyPrices_StockPrices{}
	dailyPrices.Measurements = map[string]*dailyprices_pb.Measurements{}
	for rows.Next() {
		var stockPrices dailyprices_pb.DailyPrices_StockPrices
		var ticker string
		var exchange string
		err := rows.Scan(
			&ticker,
			&stockPrices.Open,
			&stockPrices.Close,
			&stockPrices.Low,
			&stockPrices.High,
			&stockPrices.Volume,
			&exchange,
		)
		if err != nil {
			return nil, fmt.Errorf("Error while parsing row: `%s`", err.Error())
		}
		dailyPrices.StockPrices[ticker] = &stockPrices
		dailyPrices.Measurements[ticker] = &dailyprices_pb.Measurements{Exchange: exchange}
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("Error constructing response: `%s`", rows.Err().Error())
	}

	// Update averages.
	if err = s.updateAverages(tickTime, &dailyPrices); err != nil {
		return nil, fmt.Errorf("Error updating averages: %s", err.Error())
	}

	// Stamp, cache, and return.
	dailyPrices.Timestamp = req.GetTimestamp()
	dailyPrices.Version = req.GetVersion()
	if _, ok := s.cache[dailyPrices.GetVersion()]; !ok {
		s.cache[dailyPrices.GetVersion()] = map[time.Time]*dailyprices_pb.DailyPrices{}
	}
	s.cache[dailyPrices.GetVersion()][tickTime] = &dailyPrices
	return &dailyPrices, nil
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

	// Reset averages.
	s.rollingAverages = map[string]*mean.Rolling{}
	s.rollingAverageReturns = map[string]*mean.Rolling{}

	// Reset first seen times.
	s.firstSeen = map[string]time.Time{}
	s.lastSeen = map[string]time.Time{}
	s.missingDates = map[string]map[time.Time]struct{}{}

	// Reset alpha
	s.alpha = map[string]*alpha.Rolling{}

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
