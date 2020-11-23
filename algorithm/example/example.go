package buyspy

import (
	"context"
	"fmt"

	algorithmio_pb "github.com/d-sparks/gravy/algorithm/proto"
	"github.com/d-sparks/gravy/registrar"
	supervisor_pb "github.com/d-sparks/gravy/supervisor/proto"
	"github.com/golang/protobuf/ptypes"
)

// MyAlgorithm description.
type MyAlgorithm struct {
	algorithmio_pb.UnimplementedAlgorithmServer

	// Algorithm ID (usually "buyandhold" unless multiple are running)
	id          string
	algorithmID *supervisor_pb.AlgorithmId

	registrar *registrar.R
}

// New creates a new, uninitialized BuySPY algorithm.
func New(algorithmID string) *BuySPY {
	return &BuySPY{
		id:          algorithmID,
		algorithmID: &supervisor_pb.AlgorithmId{AlgorithmId: algorithmID},
	}
}

// Init initializes the registrar. The algorithm should be listening before calling Init to avoid deadlocks.
func (b *BuySPY) Init() error {
	var err error
	b.registrar, err = registrar.NewWithSupervisor()
	return err
}

// Close closes the regitsrar.
func (b *BuySPY) Close() {
	b.registrar.Close()
}

// Execute implements the algorithm interface.
func (b *BuySPY) Execute(ctx context.Context, input *algorithmio_pb.Input) (*algorithmio_pb.Output, error) {
	fmt.Printf("Excuting algorithm on %s\n", ptypes.TimestampString(input.GetTimestamp()))

	// req := dailyprices_pb.Request{Timestamp: input.GetTimestamp(), Version: 0}
	// _, err := b.registrar.Supervisor.PlaceOrder(ctx, order); err != nil {

	if _, err := b.registrar.Supervisor.DoneTrading(ctx, b.algorithmID); err != nil {
		return nil, fmt.Errorf("Error calling DoneTrading from `%s`: %s", b.id, err.Error())
	}

	return &algorithmio_pb.Output{}, nil
}
