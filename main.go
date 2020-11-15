package main

import (
	"context"
	"log"
	"time"

	"github.com/d-sparks/gravy/registrar"
	supervisor_pb "github.com/d-sparks/gravy/supervisor/proto"
	"github.com/golang/protobuf/ptypes"
	timestamp_pb "github.com/golang/protobuf/ptypes/timestamp"
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

func main() {
	// Open registrar.
	registrar, err := registrar.NewWithSupervisor()
	fatalIfErr(err)
	defer registrar.Close()

	// Open client.
	var req supervisor_pb.SynchronousDailySimInput
	req.Start = parseTimeOrDie("2006-01-03")
	req.End = parseTimeOrDie("2006-03-03")
	req.OutputDir = "/tmp/foo"

	// Send request.
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	_, err = registrar.Supervisor.SynchronousDailySim(ctx, &req)
	fatalIfErr(err)
}
