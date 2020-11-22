package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/d-sparks/gravy/algorithm/buyandhold"
	algorithmio_pb "github.com/d-sparks/gravy/algorithm/proto"
	"google.golang.org/grpc"
)

var (
	id              = flag.String("id", "buyandhold", "Algorithm ID.")
	port            = flag.Int("port", 17502, "Port for rpc server.")
	rebalancePeriod = flag.Int("rebalance_period", 20, "Rebalance period in trading ticks/days.")
)

func main() {
	flag.Parse()

	// Listen on tcp
	listeningOn := fmt.Sprintf("localhost:%d", *port)
	lis, err := net.Listen("tcp", listeningOn)
	if err != nil {
		log.Fatalf("Failed to listen over tcp: %s", err.Error())
	}

	// Make server (uninitialized)
	algorithmServer := buyandhold.New(*id, *rebalancePeriod)

	// Create grcp server and serve
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	algorithmio_pb.RegisterAlgorithmServer(grpcServer, algorithmServer)

	// Init and serve.
	algorithmServer.Init()
	defer algorithmServer.Close()
	log.Printf("Listening on `%s`", listeningOn)
	grpcServer.Serve(lis)
}
