package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/d-sparks/gravy/data/dailyprices"
	dailyprices_pb "github.com/d-sparks/gravy/data/dailyprices/proto"
	"google.golang.org/grpc"
)

var (
	port        = flag.Int("port", 17501, "Port for rpc server.")
	postgresURL = flag.String(
		"postgres_url",
		"postgres://localhost:5432/gravy?sslmode=disable",
		"Gravy db url.",
	)
	dailyPricesTable  = flag.String("prices_table", "dailyprices", "Daily prices logical table.")
	tradingDatesTable = flag.String("trading_dates", "tradingdates", "Trading dates logical table.")
)

func main() {
	flag.Parse()

	// Listen on tcp
	listeningOn := fmt.Sprintf("localhost:%d", *port)
	lis, err := net.Listen("tcp", listeningOn)
	if err != nil {
		log.Fatalf("Failed to listen over tcp: %s", err.Error())
	}

	// Make daily prices server (connect to DB)
	dailyPricesServer, err := dailyprices.NewServer(*postgresURL, *dailyPricesTable, *tradingDatesTable)
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
