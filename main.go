package main

import (
	"context"
	"fmt"
	"log"
	"time"

	dailyprices_pb "github.com/d-sparks/gravy/data/dailyprices/proto"
	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc"
)

func main() {
	// Open client.
	opts := []grpc.DialOption{grpc.WithInsecure(), grpc.WithBlock()}
	conn, err := grpc.Dial("localhost:17501", opts...)
	if err != nil {
		log.Fatalf(err.Error())
	}
	defer conn.Close()
	pricesClient := dailyprices_pb.NewDataClient(conn)

	// Construct Request.
	var req dailyprices_pb.Request
	date, err := time.Parse("2006-01-02", "2006-01-03")
	if err != nil {
		log.Fatalf(err.Error())
	}
	timestamp, err := ptypes.TimestampProto(date)
	if err != nil {
		log.Fatalf(err.Error())
	}
	req.Timestamp = timestamp

	// Send request.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	prices, err := pricesClient.Get(ctx, &req)
	if err != nil {
		log.Fatalf(err.Error())
	}

	// Print result.
	for _, stockPrices := range prices.GetStockPrices() {
		fmt.Println(stockPrices)
		break
	}

	// Trading dates request.
	var dateRange dailyprices_pb.Range
	dateRange.Lb = timestamp
	dateRange.Ub = timestamp

	ctxx, cancell := context.WithTimeout(context.Background(), time.Second)
	defer cancell()
	tradingDates, err := pricesClient.TradingDatesInRange(ctxx, &dateRange)
	if err != nil {
		log.Fatalf(err.Error())
	}

	for _, date := range tradingDates.Timestamps {
		fmt.Println(ptypes.TimestampString(date))
	}
}
