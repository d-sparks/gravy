package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/d-sparks/gravy/supervisor"
	supervisor_pb "github.com/d-sparks/gravy/supervisor/proto"
	"google.golang.org/grpc"
)

var (
	port         = flag.Int("port", 17500, "Port for rpc server.")
	mode         = flag.String("mode", "sync", "Supervision mode (sync, async, paper, live).")
	timescaleURL = flag.String(
		"timescale_url",
		"postgres://localhost:5432/gravy_timescale_output?sslmode=disable",
		"Timescale DB url",
	)
)

func parseTradingMode(mode string) supervisor.TradingMode {
	switch mode {
	case "sync":
		return supervisor.SyncTM
	case "async":
		return supervisor.AsyncTM
	case "paper":
		return supervisor.PaperTM
	case "live":
		return supervisor.LiveTM
	}
	log.Fatalf("Unrecognized mode: %s", mode)
	return supervisor.SyncTM
}

func main() {
	flag.Parse()

	// Listen on tcp
	listeningOn := fmt.Sprintf("localhost:%d", *port)
	lis, err := net.Listen("tcp", listeningOn)
	if err != nil {
		log.Fatalf("Failed to listen over tcp: %s", err.Error())
	}

	// Make server (uninitialized)
	tradingMode := parseTradingMode(*mode)
	supervisorServer, err := supervisor.New(tradingMode, *timescaleURL)
	if err != nil {
		log.Fatalf("Error constructing server: %s", err.Error())
	}

	// Create grcp server and serve
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	supervisor_pb.RegisterSupervisorServer(grpcServer, supervisorServer)

	// Init and serve.
	timescaleID, err := supervisorServer.Init()
	defer supervisorServer.Close()
	if err != nil {
		log.Fatalf("Error initializing supervisor: `%s`", err.Error())
	}
	log.Printf("Listening on `%s`", listeningOn)
	log.Printf("Writing to timescale table `%s`", timescaleID)
	grpcServer.Serve(lis)
}
