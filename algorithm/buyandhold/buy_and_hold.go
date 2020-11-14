package buyandhold

import (
	"context"
	"fmt"
	"math"

	buyandhold_pb "github.com/d-sparks/gravy/algorithm/buyandhold/proto"
	algorithmio_pb "github.com/d-sparks/gravy/algorithm/proto"
	dailyprices_pb "github.com/d-sparks/gravy/data/dailyprices/proto"
	"github.com/d-sparks/gravy/gravy"
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

	invested bool

	registrar *registrar.R
}

// New creates a new, uninitialized BuyAndHold algorithm.
func New() *BuyAndHold {
	var b BuyAndHold

	b.algorithmID = &supervisor_pb.AlgorithmId{}
	b.algorithmID.AlgorithmId = registrar.AlgorithmSpecs[algorithmEnum].ID
	b.invested = false

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

// InvestApproximatelyUniformly attempts to invest approximately uniformly.
func (b *BuyAndHold) InvestApproximatelyUniformly(
	portfolio *supervisor_pb.Portfolio,
	prices *dailyprices_pb.DailyPrices,
) (orders []*supervisor_pb.Order) {
	totalLimitOfOrders := 0.0
	target := 0.99 * gravy.TargetUniformInvestment(portfolio, prices)
	targetInvestments := map[string]float64{}

	// Note: this will represent a total limit of
	//
	//   1.01 * \sum_stock floor(target / price) * price <=
	//   1.01 * \sum_stock target =
	//   1.01 * 0.99 * \sum_stocks portfolioValue / # stocks =
	//   1.01 * 0.99 * portfolioValue =
	//   0.9999 * portfolioValue
	//
	// Thus, the investment is safe.
	for ticker, stockPrices := range prices.GetStockPrices() {
		var order supervisor_pb.Order
		order.AlgorithmId = b.algorithmID
		order.Ticker = ticker
		order.Volume = math.Floor(target / stockPrices.GetClose())
		if order.GetVolume() == 0.0 {
			continue
		}
		order.Limit = 1.01 * stockPrices.GetClose()
		totalLimitOfOrders += order.GetVolume() * order.GetLimit()
		orders = append(orders, &order)
	}

	// Continue investing until we can't get closer to a uniform investment.
	for portfolio.GetUsd()-totalLimitOfOrders > 0 {
		var nextTicker string
		nextImprovement := 0.0
		nextPrice := 0.0
		for ticker, stockPrices := range prices.GetStockPrices() {
			currentTarget := targetInvestments[ticker]
			closePrice := stockPrices.GetClose()
			if closePrice+totalLimitOfOrders > portfolio.GetUsd() {
				continue
			}
			hypotheticalDelta := math.Abs(closePrice + currentTarget - target)
			currentDelta := math.Abs(currentTarget - target)
			improvement := currentDelta - hypotheticalDelta
			if improvement > nextImprovement {
				nextImprovement = improvement
				nextTicker = ticker
				nextPrice = closePrice
			}
		}
		if nextImprovement == 0.0 {
			// No improvement can be made.
			break
		}
		// Place an order for a nextTicker
		var order supervisor_pb.Order
		order.AlgorithmId = b.algorithmID
		order.Ticker = nextTicker
		order.Volume = 1.0
		order.Limit = 1.01 * nextPrice
		totalLimitOfOrders += order.GetVolume() * order.GetLimit()
		orders = append(orders, &order)
	}

	return
}

// Execute implements the algorithm interface.
func (b *BuyAndHold) Execute(ctx context.Context, input *algorithmio_pb.Input) (*algorithmio_pb.Output, error) {
	fmt.Printf("Excuting algorithm on %s\n", ptypes.TimestampString(input.GetTimestamp()))

	portfolio, err := b.registrar.Supervisor.GetPortfolio(ctx, b.algorithmID)
	if err != nil {
		return nil, fmt.Errorf(
			"Error getting portfolio in `%s`: %s",
			registrar.AlgorithmSpecs[algorithmEnum].ID,
			err.Error(),
		)
	}

	var req dailyprices_pb.Request
	req.Timestamp = input.GetTimestamp()
	req.Version = 0
	dailyPrices, err := b.registrar.DailyPrices.Get(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf(
			"Error getting daily prices in `%s`: %s",
			registrar.AlgorithmSpecs[algorithmEnum].ID,
			err.Error(),
		)
	}

	if !b.invested {
		orders := b.InvestApproximatelyUniformly(portfolio, dailyPrices)
		for _, order := range orders {
			if _, err := b.registrar.Supervisor.PlaceOrder(ctx, order); err != nil {
				return nil, fmt.Errorf(
					"Error placing supplemental order from `%s`: %s",
					registrar.AlgorithmSpecs[algorithmEnum].ID,
					err.Error(),
				)
			}
		}
		b.invested = true
	}

	if _, err = b.registrar.Supervisor.DoneTrading(ctx, b.algorithmID); err != nil {
		return nil, fmt.Errorf(
			"Error calling DoneTrading from `%s`: %s",
			registrar.AlgorithmSpecs[algorithmEnum].ID,
			err.Error())
	}

	return &algorithmio_pb.Output{}, nil
}
