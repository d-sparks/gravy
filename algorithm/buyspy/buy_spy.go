package buyspy

import (
	"context"
	"fmt"
	"math"

	algorithmio_pb "github.com/d-sparks/gravy/algorithm/proto"
	dailyprices_pb "github.com/d-sparks/gravy/data/dailyprices/proto"
	"github.com/d-sparks/gravy/registrar"
	supervisor_pb "github.com/d-sparks/gravy/supervisor/proto"
	"github.com/golang/protobuf/ptypes"
)

// BuySPY goes all in on SPY (Alpha 0, Beta 1 strategy).
type BuySPY struct {
	algorithmio_pb.UnimplementedAlgorithmServer

	// Algorithm ID (usually "buyandhold" unless multiple are running)
	id          string
	algorithmID *supervisor_pb.AlgorithmId

	invested bool

	registrar *registrar.R
}

// New creates a new, uninitialized BuySPY algorithm.
func New(algorithmID string) *BuySPY {
	return &BuySPY{
		id:          algorithmID,
		algorithmID: &supervisor_pb.AlgorithmId{AlgorithmId: algorithmID},
		invested:    false,
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

// InvestInSPY attempts to invest entirely in SPY.
func (b *BuySPY) InvestInSPY(
	portfolio *supervisor_pb.Portfolio,
	prices *dailyprices_pb.DailyPrices,
) (orders []*supervisor_pb.Order) {
	// Get prices and desired number of units.
	SPYPrice := prices.GetStockPrices()["SPY"].GetClose()
	limit := 1.01 * SPYPrice
	USD := portfolio.GetUsd()
	units := math.Floor(USD / limit)

	// Create orders.
	orders = append(orders, &supervisor_pb.Order{
		AlgorithmId: b.algorithmID, Ticker: "SPY", Volume: units, Limit: limit,
	})

	return
}

// Execute implements the algorithm interface.
func (b *BuySPY) Execute(ctx context.Context, input *algorithmio_pb.Input) (*algorithmio_pb.Output, error) {
	fmt.Printf("Excuting algorithm on %s\n", ptypes.TimestampString(input.GetTimestamp()))

	req := dailyprices_pb.Request{Timestamp: input.GetTimestamp(), Version: 0}

	if !b.invested {
		portfolio, err := b.registrar.Supervisor.GetPortfolio(ctx, b.algorithmID)
		if err != nil {
			return nil, fmt.Errorf("Error getting portfolio in `%s`: %s", b.id, err.Error())
		}

		dailyPrices, err := b.registrar.DailyPrices.Get(ctx, &req)
		if err != nil {
			return nil, fmt.Errorf("Error getting daily prices in `%s`: %s", b.id, err.Error())
		}
		orders := b.InvestInSPY(portfolio, dailyPrices)
		for _, order := range orders {
			if _, err := b.registrar.Supervisor.PlaceOrder(ctx, order); err != nil {
				return nil, fmt.Errorf(
					"Error placing order from `%s`: %s", b.id, err.Error(),
				)
			}
		}
		b.invested = true
	}

	if _, err := b.registrar.Supervisor.DoneTrading(ctx, b.algorithmID); err != nil {
		return nil, fmt.Errorf("Error calling DoneTrading from `%s`: %s", b.id, err.Error())
	}

	return &algorithmio_pb.Output{}, nil
}
