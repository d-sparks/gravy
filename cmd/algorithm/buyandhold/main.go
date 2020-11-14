package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/d-sparks/gravy/algorithm/buyandhold"
	buyandhold_pb "github.com/d-sparks/gravy/algorithm/buyandhold/proto"
	"google.golang.org/grpc"
)

var port = flag.Int("port", 17502, "Port for rpc server.")

func main() {
	flag.Parse()

	// Listen on tcp
	listeningOn := fmt.Sprintf("localhost:%d", *port)
	lis, err := net.Listen("tcp", listeningOn)
	if err != nil {
		log.Fatalf("Failed to listen over tcp: %s", err.Error())
	}

	// Make server (uninitialized)
	algorithmServer := buyandhold.New()

	// Create grcp server and serve
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	buyandhold_pb.RegisterBuyAndHoldServer(grpcServer, algorithmServer)

	// Init and serve.
	algorithmServer.Init()
	defer algorithmServer.Close()
	log.Printf("Listening on `%s`", listeningOn)
	grpcServer.Serve(lis)
}
