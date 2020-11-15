package registrar

import (
	"fmt"

	"github.com/d-sparks/gravy/algorithm"
	buyandhold_pb "github.com/d-sparks/gravy/algorithm/buyandhold/proto"
	dailyprices_pb "github.com/d-sparks/gravy/data/dailyprices/proto"
	supervisor_pb "github.com/d-sparks/gravy/supervisor/proto"
	"google.golang.org/grpc"
)

const (
	supervisorURL  = "localhost:17500"
	dailypricesURL = "localhost:17501"
)

// AlgorithmEnum is an enumeration of the algorithms known to the registrar.
type AlgorithmEnum int

const (
	// BuyAndHold algorithm.
	BuyAndHold AlgorithmEnum = iota
)

// AlgorithmSpec holds the URL and name of an algorithm.
type AlgorithmSpec struct {
	ID  string
	URL string
}

// AlgorithmSpecs holds an algorithmSpec for each enumerated algorithm.
var AlgorithmSpecs = map[AlgorithmEnum]*AlgorithmSpec{
	BuyAndHold: &AlgorithmSpec{"buyandhold", "localhost:17502"},
}

// R is the registrar. Has clients for all the grpc services comprising gravy.
type R struct {
	// Supervisor
	Supervisor supervisor_pb.SupervisorClient

	// Algorithms
	Algorithms map[string]algorithm.A

	// Data
	DailyPrices dailyprices_pb.DataClient

	connections          []*grpc.ClientConn
	algorithmConnections []*grpc.ClientConn
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

// openAlgorithm attempts to open a given algorithm.
func (r *R) openAlgorithm(algorithmEnum AlgorithmEnum) error {
	// Get spec.
	spec, ok := AlgorithmSpecs[algorithmEnum]
	if !ok {
		return fmt.Errorf("Unknown algorithm %d", algorithmEnum)
	}

	// Connect to the algorithm server.
	opts := []grpc.DialOption{grpc.WithInsecure(), grpc.WithBlock()}
	conn, err := grpc.Dial(spec.URL, opts...)
	if err != nil {
		return fmt.Errorf("Error connecting to algorithm %d: %s", algorithmEnum, err.Error())
	}
	r.algorithmConnections = append(r.algorithmConnections, conn)

	// Create the client.
	switch algorithmEnum {
	case BuyAndHold:
		r.Algorithms[spec.ID] = buyandhold_pb.NewBuyAndHoldClient(conn)
	}

	return nil
}

// InitAlgorithms initializes a set of algorithsm.
func (r *R) InitAlgorithms(algorithms ...AlgorithmEnum) error {
	r.CloseAllAlgorithms()
	r.Algorithms = map[string]algorithm.A{}
	for _, algorithmEnum := range algorithms {
		if err := r.openAlgorithm(algorithmEnum); err != nil {
			return err
		}
	}
	return nil
}

// InitAllAlgorithms initializes all known algorithms.
func (r *R) InitAllAlgorithms() error {
	r.CloseAllAlgorithms()
	r.Algorithms = map[string]algorithm.A{}
	for algorithmEnum := range AlgorithmSpecs {
		if err := r.openAlgorithm(algorithmEnum); err != nil {
			return err
		}
	}
	return nil
}

// CloseAllAlgorithms closes all open algorithms.
func (r *R) CloseAllAlgorithms() {
	for _, conn := range r.algorithmConnections {
		conn.Close()
	}
	r.algorithmConnections = []*grpc.ClientConn{}
}

// New constructs a new registrar. Does not connect to supervisor (this is called from supervisor).
func New() (*R, error) {
	r := R{}

	// Open connections to data sources.
	if err := r.openDailyPricesConnection(dailypricesURL); err != nil {
		return nil, err
	}

	return &r, nil
}

// NewWithSupervisor creates a registrar including the supervisor.
func NewWithSupervisor() (*R, error) {
	// New registrar.
	r, err := New()
	if err != nil {
		return nil, err
	}

	// Open connection to the supervisor.
	if err := r.openSupervisorConnection(supervisorURL); err != nil {
		return nil, err
	}

	return r, nil
}
