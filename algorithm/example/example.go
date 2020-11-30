package myalgorithm

import (
	"context"
	"fmt"

	algorithmio_pb "github.com/d-sparks/gravy/algorithm/proto"
	dailyprices_pb "github.com/d-sparks/gravy/data/dailyprices/proto"
	"github.com/d-sparks/gravy/registrar"
	supervisor_pb "github.com/d-sparks/gravy/supervisor/proto"
	"github.com/golang/protobuf/ptypes"
)

// MyAlgorithm is a simple algorithm that tries to buy a fairly diversified portfolio and holds. It sells everything and
// rebalances periodically.
type MyAlgorithm struct {
	algorithmio_pb.UnimplementedAlgorithmServer

	// Algorithm ID (usually "buyandhold" unless multiple are running)
	id          string
	algorithmID *supervisor_pb.AlgorithmId
	registrar   *registrar.R
}

// skipTrading is a precondition. Save time if you don't need to fetch prices/portfolio.
func (b *MyAlgorithm) skipTrading() bool {
	return false
}

// trade is the algorithm itself.
func (b *MyAlgorithm) trade(
	portfolio *supervisor_pb.Portfolio,
	data *dailyprices_pb.DailyData,
) []*supervisor_pb.Order {
	return nil
}

// New creates a new, uninitialized MyAlgorithm algorithm.
func New(algorithmID string, rebalancePeriod int) *MyAlgorithm {
	return &MyAlgorithm{
		id:          algorithmID,
		algorithmID: &supervisor_pb.AlgorithmId{AlgorithmId: algorithmID},
	}
}

// ******************************
//  Mostly boilerplate hereafter
// ******************************

// Init initializes the registrar. The algorithm should be listening before calling Init to avoid deadlocks.
func (b *MyAlgorithm) Init() error {
	var err error
	b.registrar, err = registrar.NewWithSupervisor()
	return err
}

// Close closes the regitsrar.
func (b *MyAlgorithm) Close() {
	b.registrar.Close()
}

// Execute implements the algorithm interface.
func (b *MyAlgorithm) Execute(ctx context.Context, input *algorithmio_pb.Input) (*algorithmio_pb.Output, error) {
	fmt.Printf("Excuting algorithm on %s\n", ptypes.TimestampString(input.GetTimestamp()))
	orders := []*supervisor_pb.Order{}

	if !b.skipTrading() {
		portfolio, err := b.registrar.Supervisor.GetPortfolio(ctx, b.algorithmID)
		if err != nil {
			return nil, fmt.Errorf("Error getting portfolio in `%s`: %s", b.id, err.Error())
		}

		req := dailyprices_pb.Request{Timestamp: input.GetTimestamp(), Version: 0}
		data, err := b.registrar.DailyPrices.Get(ctx, &req)
		if err != nil {
			return nil, fmt.Errorf("Error getting daily prices in `%s`: %s", b.id, err.Error())
		}

		orders = b.trade(portfolio, data)
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
