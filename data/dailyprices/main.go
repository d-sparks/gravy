package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	dailyprices_pb "github.com/d-sparks/gravy/data/dailyprices/proto"
	"github.com/golang/protobuf/ptypes"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
)

var (
	port        = flag.Int("port", 17501, "Port for rpc server.")
	postgresURL = flag.String(
		"postgres_url",
		"postgres://localhost:5432/gravy?sslmode=disable",
		"Gravy db url.",
	)
	dailyPricesTable = flag.String("prices_table", "dailyprices", "Daily prices logical table.")
)

// DailyPricesServer implements dailyprices_pb.DataServer.
type DailyPricesServer struct {
	dailyprices_pb.UnimplementedDataServer

	// PostGRES
	db        *sql.DB
	tableName string

	// Cache
	mu    sync.Mutex
	cache map[int32]map[time.Time]*dailyprices_pb.DailyPrices
}

// NewDailyPricesServer creates an empty daily prices server.
func NewDailyPricesServer(postgresURL string, dailyPricesTable string) (*DailyPricesServer, error) {
	log.Printf("Connecting to database `%s`", postgresURL)
	db, err := sql.Open("postgres", postgresURL)
	if err != nil {
		return nil, fmt.Errorf("Error connecting to postgres: %s", err.Error())
	}

	var dailyPricesServer DailyPricesServer
	dailyPricesServer.db = db
	dailyPricesServer.tableName = dailyPricesTable
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
			return cachedDailyPrices, nil
		}
	}
	s.mu.Unlock()

	// Query database.
	rows, err := s.db.Query(
		fmt.Sprintf(
			"SELECT ticker, open, close, adj_close, low, high, volume FROM %s WHERE date = $1",
			s.tableName,
		),
		tickTime.Format("2006-01-02"),
	)
	if err != nil {
		return nil, fmt.Errorf("Error reading from db: `%s`", err.Error())
	}

	// Construct daily prices by scanning the query result.
	var dailyPrices dailyprices_pb.DailyPrices
	for rows.Next() {
		var stockPrices dailyprices_pb.DailyPrices_StockPrices
		err := rows.Scan(
			&stockPrices.Ticker,
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
		dailyPrices.StockPrices = append(dailyPrices.GetStockPrices(), &stockPrices)
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

func main() {
	flag.Parse()

	// Listen on tcp
	listeningOn := fmt.Sprintf("localhost:%d", *port)
	lis, err := net.Listen("tcp", listeningOn)
	if err != nil {
		log.Fatalf("Failed to listen over tcp: %s", err.Error())
	}

	// Make daily prices server (connect to DB)
	dailyPricesServer, err := NewDailyPricesServer(*postgresURL, *dailyPricesTable)
	if err != nil {
		log.Fatalf("Error constructing server: %s", err.Error())
	}

	// Create grcp server and serve
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	dailyprices_pb.RegisterDataServer(grpcServer, dailyPricesServer)
	log.Printf("Listening on `%s`", listeningOn)
	grpcServer.Serve(lis)
}
