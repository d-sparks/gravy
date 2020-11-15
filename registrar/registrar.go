package registrar

import (
	"fmt"

	"github.com/d-sparks/gravy/algorithm"
	algorithmio_pb "github.com/d-sparks/gravy/algorithm/proto"
	dailyprices_pb "github.com/d-sparks/gravy/data/dailyprices/proto"
	supervisor_pb "github.com/d-sparks/gravy/supervisor/proto"
	"google.golang.org/grpc"
)

const (
	// TODO: These should probably also be in the input to the sim commands.
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
func (r *R) openAlgorithm(spec *supervisor_pb.AlgorithmSpec) error {
	// Connect to the algorithm server.
	fmt.Printf("Connecting to %s at %s\n", spec.GetId().GetAlgorithmId(), spec.GetUrl())
	opts := []grpc.DialOption{grpc.WithInsecure(), grpc.WithBlock()}
	conn, err := grpc.Dial(spec.GetUrl(), opts...)
	if err != nil {
		return fmt.Errorf("Error connecting to algorithm %s: %s", spec.GetId().GetAlgorithmId(), err.Error())
	}
	r.algorithmConnections = append(r.algorithmConnections, conn)

	// Create the client.
	fmt.Printf("Connected intermediately to %s at %s\n", spec.GetId().GetAlgorithmId(), spec.GetUrl())
	r.Algorithms[spec.GetId().GetAlgorithmId()] = algorithmio_pb.NewAlgorithmClient(conn)
	fmt.Printf("Connected to %s at %s\n", spec.GetId().GetAlgorithmId(), spec.GetUrl())

	return nil
}

// InitAlgorithms initializes a set of algorithsm.
func (r *R) InitAlgorithms(algorithms ...*supervisor_pb.AlgorithmSpec) error {
	r.CloseAllAlgorithms()
	r.Algorithms = map[string]algorithm.A{}
	for _, algorithmSpec := range algorithms {
		if err := r.openAlgorithm(algorithmSpec); err != nil {
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
