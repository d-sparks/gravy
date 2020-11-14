package buyandhold

import (
	"context"
	"fmt"

	buyandhold_pb "github.com/d-sparks/gravy/algorithm/buyandhold/proto"
	algorithmio_pb "github.com/d-sparks/gravy/algorithm/proto"
	"github.com/d-sparks/gravy/registrar"
	supervisor_pb "github.com/d-sparks/gravy/supervisor/proto"
	"github.com/golang/protobuf/ptypes"
)

const algorithmEnum = registrar.BuyAndHold

// BuyAndHold is a simple algorithm that tries to buy a fairly diversified portfolio and holds forever. If stocks are
// delisted, the proceeds are invested in an attempt to extend diversity.
type BuyAndHold struct {
	buyandhold_pb.UnimplementedBuyAndHoldServer

	algorithmID *supervisor_pb.AlgorithmId

	registrar *registrar.R
}

// New creates a new, uninitialized BuyAndHold algorithm.
func New() *BuyAndHold {
	var b BuyAndHold

	b.algorithmID = &supervisor_pb.AlgorithmId{}
	b.algorithmID.AlgorithmId = registrar.AlgorithmSpecs[algorithmEnum].ID

	return &b
}

// Init initializes the registrar. The algorithm should be listening before calling Init to avoid deadlocks.
func (b *BuyAndHold) Init() error {
	var err error
	b.registrar, err = registrar.NewWithSupervisor()
	return err
}

// Close closes the regitsrar.
func (b *BuyAndHold) Close() {
	b.registrar.Close()
}

// Execute implements the algorithm interface.
func (b *BuyAndHold) Execute(ctx context.Context, input *algorithmio_pb.Input) (*algorithmio_pb.Output, error) {
	fmt.Printf("Excuting algorithm on %s\n", ptypes.TimestampString(input.GetTimestamp()))

	_, err := b.registrar.Supervisor.DoneTrading(ctx, b.algorithmID)
	if err != nil {
		return nil, fmt.Errorf("Error telling supervisor I'm done: %s", err.Error())
	}

	return &algorithmio_pb.Output{}, nil
}
