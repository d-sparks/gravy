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
	"github.com/d-sparks/gravy/data/movingaverage"
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
	i         int

	// Store 15, 35, and 252 day rolling averages. (~20, 50, 365 days worth of trading days.)
	rollingAverages map[int]map[string]*movingaverage.Rolling

	// 252 day rolling average from 252 days ago.
	oldAverages map[string]*movingaverage.Rolling

	// First time at which the given stock was seen.
	firstSeen map[string]time.Time
	lastSeen  map[string]time.Time

	// Alpha for each ticker.
	alpha map[string]*alpha.Rolling
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

	var dailyPricesServer Server
	dailyPricesServer.db = db
	dailyPricesServer.pricesTableName = dailyPricesTable
	dailyPricesServer.tradingDatesTableName = tradingDatesTable
	dailyPricesServer.cache = map[int32]map[time.Time]*dailyprices_pb.DailyPrices{}

	return &dailyPricesServer, nil
}

// Close the underlying postgres connection.
func (s *Server) Close() {
	s.db.Close()
}

// updateAveragesForSymbol updates for a single tick.
func (s *Server) updateAveragesForSymbol(tickTime time.Time, ticker string, dailyPrices *dailyprices_pb.DailyPrices) {
	tickIx := s.timeIndex[tickTime]
	price := dailyPrices.GetStockPrices()[ticker].GetClose()

	// Get old prices.
	oldPrices := map[int]float64{}
	for _, daysAgo := range []int{15, 35, 252, 504} {
		oldPrices[daysAgo] = s.cache[0][s.times[tickIx-daysAgo]].GetStockPrices()[ticker].GetClose()
	}

	// Update oldAverages if we've observed this stock for 252 days.
	if observations := tickIx - s.timeIndex[s.firstSeen[ticker]]; observations >= 252 {
		// On the 252 day, start tracking the year old moving averages.
		if observations == 252 {
			s.oldAverages[ticker] = movingaverage.NewRolling(252)
		}
		s.oldAverages[ticker].Observe(oldPrices[252], oldPrices[504])
	}

	// Update rolling averages.
	s.rollingAverages[15][ticker].Observe(price, oldPrices[15])
	s.rollingAverages[35][ticker].Observe(price, oldPrices[35])
	s.rollingAverages[252][ticker].Observe(price, oldPrices[252])

	// Record the output
	dailyPrices.Measurements[ticker] = &dailyprices_pb.Measurements{}
	dailyPrices.Measurements[ticker].MovingAverages[15] = s.rollingAverages[15][ticker].Value()
	dailyPrices.Measurements[ticker].MovingAverages[35] = s.rollingAverages[35][ticker].Value()
	dailyPrices.Measurements[ticker].MovingAverages[252] = s.rollingAverages[252][ticker].Value()
}

// updateAlphasForSymbol updates for a single tick.
func (s *Server) updateAlphasForSymbol(
	tickTime time.Time,
	ticker string,
	spy float64,
	spy0 float64,
	spyMu float64,
	spyMu0 float64,
	dailyPrices *dailyprices_pb.DailyPrices,
) {
	// Get asset data.
	x := dailyPrices.GetStockPrices()[ticker].GetClose()
	mu := s.rollingAverages[252][ticker].Value()
	x0 := 0.0
	mu0 := 0.0
	tickIx := s.timeIndex[tickTime]
	if observations := tickIx - s.timeIndex[s.firstSeen[ticker]]; observations >= 252 {
		x0 = s.cache[0][s.times[tickIx-252]].GetStockPrices()[ticker].GetClose()
		mu0 = s.oldAverages[ticker].Value()
	}

	// Update alphas.
	s.alpha[ticker].Observe(x, spy, mu, spyMu, x0, spy0, mu0, spyMu0)

	// Return values.
	dailyPrices.Measurements[ticker].Alpha = s.alpha[ticker].Alpha()
	dailyPrices.Measurements[ticker].Beta = s.alpha[ticker].Beta()
}

// updateAverages updates the various tracked rolling averages at the given time. Also updates the given prices pointer.
func (s *Server) updateAverages(tickTime time.Time, dailyPrices *dailyprices_pb.DailyPrices) {
	// Update firstSeen if this is the first time seeing this ticker.
	for ticker := range dailyPrices.GetStockPrices() {
		if _, ok := s.firstSeen[ticker]; !ok {
			s.firstSeen[ticker] = tickTime
			s.rollingAverages[15][ticker] = movingaverage.NewRolling(15)
			s.rollingAverages[35][ticker] = movingaverage.NewRolling(35)
			s.rollingAverages[252][ticker] = movingaverage.NewRolling(252)
		}
	}

	// Update averages.
	for ticker := range s.firstSeen {
		s.updateAveragesForSymbol(tickTime, ticker, dailyPrices)
	}

	// Get SPY data for benchmark.
	spy := dailyPrices.GetStockPrices()["SPY"].GetClose()
	spyMu := s.rollingAverages[252]["SPY"].Value()
	spy0 := 0.0
	spyMu0 := 0.0
	if oldIx := s.timeIndex[tickTime] - 252; oldIx >= 0 {
		spy0 = s.cache[0][s.times[oldIx]].GetStockPrices()["SPY"].GetClose()
		spyMu0 = s.oldAverages["SPY"].Value()
	}

	// Calculate alphas.
	for ticker := range s.firstSeen {
		s.updateAlphasForSymbol(tickTime, ticker, spy, spy0, spyMu, spyMu0, dailyPrices)
	}
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
			s.mu.Unlock()
			return cachedDailyPrices, nil
		}
	}

	// Query database.
	rows, err := s.db.Query(
		fmt.Sprintf(
			"SELECT ticker, open, close, low, high, volume FROM %s WHERE date = $1",
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
	for rows.Next() {
		var stockPrices dailyprices_pb.DailyPrices_StockPrices
		var ticker string
		err := rows.Scan(
			&ticker,
			&stockPrices.Open,
			&stockPrices.Close,
			&stockPrices.Low,
			&stockPrices.High,
			&stockPrices.Volume,
		)
		if err != nil {
			return nil, fmt.Errorf("Error while parsing row: `%s`", err.Error())
		}
		dailyPrices.StockPrices[ticker] = &stockPrices
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("Error constructing response: `%s`", rows.Err().Error())
	}

	// Update averages.
	s.updateAverages(tickTime, &dailyPrices)

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
	s.i = 0

	// Get trading dates for entire session.
	dates, err := s.TradingDatesInRange(ctx, req.GetSimRange())
	if err != nil {
		return nil, fmt.Errorf("Error getting dates to start session: %s", err.Error())
	}
	s.times = make([]time.Time, len(dates.GetTimestamps()))
	for i, date := range dates.GetTimestamps() {
		if s.times[i], err = ptypes.Timestamp(date); err != nil {
			return nil,
				fmt.Errorf("Invalid date/timestamp: %s", err.Error())
		}
		s.timeIndex[s.times[i]] = i
	}

	// Reset averages.
	s.rollingAverages = map[int]map[string]*movingaverage.Rolling{
		15:  map[string]*movingaverage.Rolling{},
		35:  map[string]*movingaverage.Rolling{},
		252: map[string]*movingaverage.Rolling{},
	}
	s.oldAverages = map[string]*movingaverage.Rolling{}

	// Reset first seen times.
	s.firstSeen = map[string]time.Time{}

	// Reset alpha
	s.alpha = map[string]*alpha.Rolling{}

	return nil, nil
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
