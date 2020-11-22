package main

import (
	"context"
	"flag"
	"log"
	"strings"
	"time"

	"github.com/d-sparks/gravy/registrar"
	supervisor_pb "github.com/d-sparks/gravy/supervisor/proto"
	"github.com/golang/protobuf/ptypes"
	timestamp_pb "github.com/golang/protobuf/ptypes/timestamp"
)

var (
	start      = flag.String("start", "2005-02-25", "Start date.")
	end        = flag.String("end", "2006-02-25", "Start date.")
	outputDir  = flag.String("output_dir", "/tmp/foo", "Output directory.")
	algorithms = flag.String("algorithms", "buyandhold@localhost:17502", "Comma separated list of alg@url")
)

func fatalIfErr(err error) {
	if err != nil {
		log.Fatalf(err.Error())
	}
}

func parseTimeOrDie(timeString string) *timestamp_pb.Timestamp {
	nativeTime, err := time.Parse("2006-01-02", timeString)
	fatalIfErr(err)
	timestamp, err := ptypes.TimestampProto(nativeTime)
	fatalIfErr(err)
	return timestamp
}

func parseAlgorithmSpecOrDie(algorithmsString string) (specs []*supervisor_pb.AlgorithmSpec) {
	algorithms := strings.Split(algorithmsString, ",")
	for _, algorithm := range algorithms {
		idURL := strings.Split(algorithm, "@")
		if len(idURL) != 2 {
			log.Fatalf("Invalid spec: %s", algorithm)
		}
		spec := supervisor_pb.AlgorithmSpec{}
		spec.Id = &supervisor_pb.AlgorithmId{}
		spec.Id.AlgorithmId, spec.Url = idURL[0], idURL[1]
		if spec.GetId().GetAlgorithmId() == "" || spec.GetUrl() == "" {
			log.Fatalf("Invalid algorithm spec: %s", algorithm)
		}
		specs = append(specs, &spec)
	}
	return
}

func main() {
	flag.Parse()

	// Open registrar.
	registrar, err := registrar.NewWithSupervisor()
	fatalIfErr(err)
	defer registrar.Close()

	// Make request.
	var req supervisor_pb.SynchronousDailySimInput
	req.Start = parseTimeOrDie(*start)
	req.End = parseTimeOrDie(*end)
	req.OutputDir = *outputDir
	req.Algorithms = parseAlgorithmSpecOrDie(*algorithms)

	// Start supervisor and send request.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()
	_, err = registrar.Supervisor.SynchronousDailySim(ctx, &req)
	fatalIfErr(err)
}
