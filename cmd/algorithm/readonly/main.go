package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	algorithmio_pb "github.com/d-sparks/gravy/algorithm/proto"
	"github.com/d-sparks/gravy/algorithm/readonly"
	"google.golang.org/grpc"
)

var (
	id   = flag.String("id", "readonly", "Algorithm ID.")
	port = flag.Int("port", 17506, "Port for rpc server.")
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
	algorithmServer := readonly.New(*id)

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
