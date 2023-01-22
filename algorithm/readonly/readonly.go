package readonly

import (
	"context"
	"fmt"

	algorithmio_pb "github.com/d-sparks/gravy/algorithm/proto"
	dailyprices_pb "github.com/d-sparks/gravy/data/dailyprices/proto"
	"github.com/d-sparks/gravy/registrar"
	supervisor_pb "github.com/d-sparks/gravy/supervisor/proto"
	"github.com/golang/protobuf/ptypes"
)

// ReadOnly only reads data/prices and doesn't trade. Used for warming up the dbs.
type ReadOnly struct {
	algorithmio_pb.UnimplementedAlgorithmServer

	// Algorithm ID (usually "readonly" unless multiple are running)
	id          string
	algorithmID *supervisor_pb.AlgorithmId
	registrar   *registrar.R
}

// skipTrading is a precondition. Save time if you don't need to fetch prices/portfolio.
func (b *ReadOnly) skipTrading() bool {
	return false
}

// trade is the algorithm itself.
func (b *ReadOnly) trade(
	portfolio *supervisor_pb.Portfolio,
	data *dailyprices_pb.DailyData,
) []*supervisor_pb.Order {
	return []*supervisor_pb.Order{}
}

// New creates a new, uninitialized ReadOnly algorithm.
func New(algorithmID string) *ReadOnly {
	return &ReadOnly{
		id:          algorithmID,
		algorithmID: &supervisor_pb.AlgorithmId{AlgorithmId: algorithmID},
	}
}

// ******************************
//  Mostly boilerplate hereafter
// ******************************

// Init initializes the registrar. The algorithm should be listening before calling Init to avoid deadlocks.
func (b *ReadOnly) Init() error {
	var err error
	b.registrar, err = registrar.NewWithSupervisor()
	return err
}

// Close closes the regitsrar.
func (b *ReadOnly) Close() {
	b.registrar.Close()
}

// Execute implements the algorithm interface.
func (b *ReadOnly) Execute(ctx context.Context, input *algorithmio_pb.Input) (*algorithmio_pb.Output, error) {
	fmt.Printf("Excuting algorithm on %s\n", ptypes.TimestampString(input.GetTimestamp()))
	orders := []*supervisor_pb.Order{}

	if !b.skipTrading() {
		portfolio, err := b.registrar.Supervisor.GetPortfolio(ctx, b.algorithmID)
		if err != nil {
			return nil, fmt.Errorf("Error getting portfolio in `%s`: %s", b.id, err.Error())
		}

		req := dailyprices_pb.Request{Timestamp: input.GetTimestamp(), Version: 0}
		prices, err := b.registrar.DailyPrices.Get(ctx, &req)
		if err != nil {
			return nil, fmt.Errorf("Error getting daily prices in `%s`: %s", b.id, err.Error())
		}

		orders = b.trade(portfolio, prices)
		for _, order := range orders {
			if _, err := b.registrar.Supervisor.PlaceOrder(ctx, order); err != nil {
				return nil, fmt.Errorf(
					"Error placing order from `%s`: %s", b.id, err.Error(),
				)
			}
		}
	}
	if _, err := b.registrar.Supervisor.DoneTrading(ctx, b.algorithmID); err != nil {
		return nil, fmt.Errorf("Error calling DoneTrading from `%s`: %s", b.id, err.Error())
	}

	return &algorithmio_pb.Output{}, nil
}
