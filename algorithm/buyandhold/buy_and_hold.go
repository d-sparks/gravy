package buyandhold

import (
	"context"
	"fmt"
	"math"

	algorithmio_pb "github.com/d-sparks/gravy/algorithm/proto"
	dailyprices_pb "github.com/d-sparks/gravy/data/dailyprices/proto"
	"github.com/d-sparks/gravy/gravy"
	"github.com/d-sparks/gravy/registrar"
	supervisor_pb "github.com/d-sparks/gravy/supervisor/proto"
	"github.com/golang/protobuf/ptypes"
)

// BuyAndHold is a simple algorithm that tries to buy a fairly diversified portfolio and holds. It sells everything and
// rebalances periodically.
type BuyAndHold struct {
	algorithmio_pb.UnimplementedAlgorithmServer

	// Algorithm ID (usually "buyandhold" unless multiple are running)
	id          string
	algorithmID *supervisor_pb.AlgorithmId

	invested bool

	nextRebalance   int
	rebalancePeriod int

	registrar *registrar.R
}

// New creates a new, uninitialized BuyAndHold algorithm.
func New(algorithmID string, rebalancePeriod int) *BuyAndHold {
	var b BuyAndHold

	b.id = algorithmID
	b.algorithmID = &supervisor_pb.AlgorithmId{AlgorithmId: b.id}
	b.invested = false
	b.rebalancePeriod = rebalancePeriod

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
		volume := math.Floor(target / stockPrices.GetClose())
		if volume == 0.0 {
			continue
		}
		limit := 1.01 * stockPrices.GetClose()
		orders = append(orders, &supervisor_pb.Order{
			AlgorithmId: b.algorithmID, Ticker: ticker, Volume: volume, Limit: limit,
		})
		targetInvestments[ticker] = volume * limit
		totalLimitOfOrders += volume * limit
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
		limit := 1.01 * nextPrice
		orders = append(orders, &supervisor_pb.Order{
			AlgorithmId: b.algorithmID, Ticker: nextTicker, Volume: 1.0, Limit: limit,
		})
		targetInvestments[nextTicker] += limit
		totalLimitOfOrders += 1.0 * limit
	}

	return
}

// Execute implements the algorithm interface.
func (b *BuyAndHold) Execute(ctx context.Context, input *algorithmio_pb.Input) (*algorithmio_pb.Output, error) {
	fmt.Printf("Excuting algorithm on %s\n", ptypes.TimestampString(input.GetTimestamp()))

	if !b.invested {
		portfolio, err := b.registrar.Supervisor.GetPortfolio(ctx, b.algorithmID)
		if err != nil {
			return nil, fmt.Errorf("Error getting portfolio in `%s`: %s", b.id, err.Error())
		}

		req := dailyprices_pb.Request{Timestamp: input.GetTimestamp(), Version: 0}
		dailyPrices, err := b.registrar.DailyPrices.Get(ctx, &req)
		if err != nil {
			return nil, fmt.Errorf("Error getting daily prices in `%s`: %s", b.id, err.Error())
		}

		orders := b.InvestApproximatelyUniformly(portfolio, dailyPrices)
		for _, order := range orders {
			if _, err := b.registrar.Supervisor.PlaceOrder(ctx, order); err != nil {
				return nil, fmt.Errorf(
					"Error placing order from `%s`: %s", b.id, err.Error(),
				)
			}
		}
		b.invested = true
		b.nextRebalance = b.rebalancePeriod
	} else if b.nextRebalance == 0 {
		portfolio, err := b.registrar.Supervisor.GetPortfolio(ctx, b.algorithmID)
		if err != nil {
			return nil, fmt.Errorf("Error getting portfolio in `%s`: %s", b.id, err.Error())
		}

		orders := gravy.SellEverythingMarketOrder(b.algorithmID, portfolio)
		for _, order := range orders {
			if _, err := b.registrar.Supervisor.PlaceOrder(ctx, order); err != nil {
				return nil, fmt.Errorf(
					"Error placing order from `%s`: %s", b.id, err.Error(),
				)
			}
		}
		b.invested = false
	} else {
		b.nextRebalance--
	}

	if _, err := b.registrar.Supervisor.DoneTrading(ctx, b.algorithmID); err != nil {
		return nil, fmt.Errorf("Error calling DoneTrading from `%s`: %s", b.id, err.Error())
	}

	return &algorithmio_pb.Output{}, nil
}
