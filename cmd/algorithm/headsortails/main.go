package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/d-sparks/gravy/algorithm/headsortails"
	algorithmio_pb "github.com/d-sparks/gravy/algorithm/proto"
	"github.com/d-sparks/gravy/gravy"
	"google.golang.org/grpc"
)

var (
	id   = flag.String("id", "headsortails", "Algorithm ID.")
	port = flag.Int("port", 17505, "Port for rpc server.")

	mode          = flag.String("mode", "inference", "Mode: `train` or `inference`.")
	samplingRatio = flag.Float64("sample_ratio", 0.1, "Ratio of examples to write in training mode.")
	filename      = flag.String("output", "/tmp/fizz/headsortails_data.csv", "Output filename for training data.")
)

func new() (*os.File, *headsortails.HeadsOrTails) {
	if *mode == "inference" {
		return nil, headsortails.New(*id)
	} else if *mode == "train" {
		f, err := os.Create(*filename)
		if err != nil {
			log.Fatalf(err.Error())
		}
		return f, headsortails.NewForTraining(*id, *samplingRatio, gravy.TimePIDSeed(), bufio.NewWriter(f))
	}
	log.Fatalf("Unrecognized mode `%s`", *mode)
	return nil, nil
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
	f, algorithmServer := new()
	if f != nil {
		defer f.Close()
	}

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
