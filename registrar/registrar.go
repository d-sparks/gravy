package registrar

import (
	"github.com/d-sparks/gravy/algorithm"
	dailyprices_pb "github.com/d-sparks/gravy/data/dailyprices/proto"
	supervisor_pb "github.com/d-sparks/gravy/supervisor/proto"
	"google.golang.org/grpc"
)

const (
	supervisorURL  = "localhost:17500"
	dailypricesURL = "localhost:17501"
)

// R is the registrar. Has clients for all the grpc services comprising gravy.
type R struct {
	// Supervisor
	Supervisor supervisor_pb.SupervisorClient

	// Algorithms
	Algorithms map[string]algorithm.A

	// Data
	DailyPrices dailyprices_pb.DataClient

	connections []*grpc.ClientConn
}

// Close all connections.
func (r *R) Close() {
	for _, conn := range r.connections {
		conn.Close()
	}
}

// openSupervisorConnection opens a connection to the supervisor
func (r *R) openSupervisorConnection(url string) error {
	opts := []grpc.DialOption{grpc.WithInsecure(), grpc.WithBlock()}
	conn, err := grpc.Dial("localhost:17500", opts...)
	if err != nil {
		return err
	}
	r.connections = append(r.connections, conn)
	r.Supervisor = supervisor_pb.NewSupervisorClient(conn)
	return nil
}

// openDailyPricesConnection opens a connection to the daily prices data source.
func (r *R) openDailyPricesConnection(url string) error {
	opts := []grpc.DialOption{grpc.WithInsecure(), grpc.WithBlock()}
	conn, err := grpc.Dial("localhost:17501", opts...)
	if err != nil {
		return err
	}
	r.DailyPrices = dailyprices_pb.NewDataClient(conn)
	return nil
}

// New constructs a new registrar.
func New() (*R, error) {
	r := R{}

	// Open connection to the supervisor.
	if err := r.openSupervisorConnection(supervisorURL); err != nil {
		return nil, err
	}

	// Open connections to data sources.
	if err := r.openDailyPricesConnection(dailypricesURL); err != nil {
		return nil, err
	}

	// Instantiate algorithms.
	// TODO :)

	return &r, nil
}
