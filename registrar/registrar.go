package registrar

import (
	"github.com/d-sparks/gravy/algorithm"
	dailyprices_pb "github.com/d-sparks/gravy/data/dailyprices/proto"
	"google.golang.org/grpc"
)

const dailypricesURL string = "localhost:17501"

// R is the registrar. Has clients for all the grpc services comprising gravy.
type R struct {
	// Algorithms
	Algorithms map[string]algorithm.A

	// Data
	DailyPrices *dailyprices_pb.DataClient
}

// openDailyPricesConnection opens a connection to the daily prices data source.
func (r *R) openDailyPricesConnection(url string) error {
	opts := []grpc.DialOption{grpc.WithInsecure(), grpc.WithBlock()}
	conn, err := grpc.Dial("localhost:17501", opts...)
	if err != nil {
		return err
	}
	dataClient := dailyprices_pb.NewDataClient(conn)
	r.DailyPrices = &dataClient
	return nil
}

// New constructs a new registrar.
func New() (*R, error) {
	r := R{}

	// Open connections to data sources.
	if err := r.openDailyPricesConnection(dailypricesURL); err != nil {
		return nil, err
	}

	// Instantiate algorithms.
	// TODO :)

	return &r, nil
}
