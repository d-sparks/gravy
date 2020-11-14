package dailyprices

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"

	dailyprices_pb "github.com/d-sparks/gravy/data/dailyprices/proto"
	"github.com/golang/protobuf/ptypes"
	_ "github.com/lib/pq"
)

// DailyPricesServer implements dailyprices_pb.DataServer.
type DailyPricesServer struct {
	dailyprices_pb.UnimplementedDataServer

	// PostGRES
	db                    *sql.DB
	pricesTableName       string
	tradingDatesTableName string

	// Cache
	mu    sync.Mutex
	cache map[int32]map[time.Time]*dailyprices_pb.DailyPrices
}

// NewDailyPricesServer creates an empty daily prices server.
func NewDailyPricesServer(
	postgresURL string,
	dailyPricesTable string,
	tradingDatesTable string,
) (*DailyPricesServer, error) {
	log.Printf("Connecting to database `%s`", postgresURL)
	db, err := sql.Open("postgres", postgresURL)
	if err != nil {
		return nil, fmt.Errorf("Error connecting to postgres: %s", err.Error())
	}

	var dailyPricesServer DailyPricesServer
	dailyPricesServer.db = db
	dailyPricesServer.pricesTableName = dailyPricesTable
	dailyPricesServer.tradingDatesTableName = tradingDatesTable
	dailyPricesServer.cache = map[int32]map[time.Time]*dailyprices_pb.DailyPrices{}

	return &dailyPricesServer, nil
}

// Close the underlying postgres connection.
func (s *DailyPricesServer) Close() {
	s.db.Close()
}

// Get implements the get endpoint for dailyprices_pb.DataServer.
func (s *DailyPricesServer) Get(ctx context.Context, req *dailyprices_pb.Request) (*dailyprices_pb.DailyPrices, error) {
	// Parse timestamp to Golang native time.
	tickTime, err := ptypes.Timestamp(req.GetTimestamp())
	if err != nil {
		return nil, fmt.Errorf("Invalid timestamp: %s", err.Error())
	}

	// Check cache.
	s.mu.Lock()
	if versionPrices, ok := s.cache[req.GetVersion()]; ok {
		if cachedDailyPrices, ok := versionPrices[tickTime]; ok {
			s.mu.Unlock()
			return cachedDailyPrices, nil
		}
	}
	s.mu.Unlock()

	// Query database.
	rows, err := s.db.Query(
		fmt.Sprintf(
			"SELECT ticker, open, close, adj_close, low, high, volume FROM %s WHERE date = $1",
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
			&stockPrices.AdjClose,
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

	// Stamp, cache, and return.
	dailyPrices.Timestamp = req.GetTimestamp()
	dailyPrices.Version = req.GetVersion()
	s.mu.Lock()
	if _, ok := s.cache[dailyPrices.GetVersion()]; !ok {
		s.cache[dailyPrices.GetVersion()] = map[time.Time]*dailyprices_pb.DailyPrices{}
	}
	s.cache[dailyPrices.GetVersion()][tickTime] = &dailyPrices
	s.mu.Unlock()
	return &dailyPrices, nil
}

// TradingDatesInRange implements the interface method. Returns trading dates in the given range.
func (s *DailyPricesServer) TradingDatesInRange(
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
