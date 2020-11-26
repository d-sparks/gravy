package buygoodandhold

import (
	"context"
	"fmt"

	algorithmio_pb "github.com/d-sparks/gravy/algorithm/proto"
	dailyprices_pb "github.com/d-sparks/gravy/data/dailyprices/proto"
	"github.com/d-sparks/gravy/gravy"
	"github.com/d-sparks/gravy/registrar"
	supervisor_pb "github.com/d-sparks/gravy/supervisor/proto"
	"github.com/golang/protobuf/ptypes"
)

// BuyGoodAndHold is a simple algorithm that tries to buy a fairly diversified portfolio and holds. It sells everything and
// rebalances periodically.
type BuyGoodAndHold struct {
	algorithmio_pb.UnimplementedAlgorithmServer

	// Algorithm ID (usually "buyandhold" unless multiple are running)
	id          string
	algorithmID *supervisor_pb.AlgorithmId
	registrar   *registrar.R

	// For business logic
	invested        bool
	nextRebalance   int
	rebalancePeriod int
}

// skipTrading is a precondition. Save time if you don't need to fetch prices/portfolio.
func (b *BuyGoodAndHold) skipTrading() bool {
	if b.invested && b.nextRebalance > 0 {
		b.nextRebalance--
		return true
	}
	return false
}

// getGoodStocks picks n "good" stocks by some ad hoc method.
func getGoodStocks(dailyPrices *dailyprices_pb.DailyPrices) map[string]struct{} {
	// Only consider stocks
	admissible := []string{}
	for ticker, prices := range dailyPrices.GetStockPrices() {
		exchange := dailyPrices.GetMeasurements()[ticker].GetExchange()
		price := prices.GetClose()
		volume := prices.GetVolume()
		if price >= 5.0 && price*volume >= float64(2E7) && (exchange == "NYSE" || exchange == "NASDAQ") {
			admissible = append(admissible, ticker)
		}
	}

	// Pick the first n admissible (lowest beta).
	good := map[string]struct{}{}
	for i := 0; i < len(admissible); i++ {
		good[admissible[i]] = struct{}{}
	}

	return good
}

// trade is the algorithm itself.
func (b *BuyGoodAndHold) trade(
	portfolio *supervisor_pb.Portfolio,
	dailyPrices *dailyprices_pb.DailyPrices,
) []*supervisor_pb.Order {
	if !b.invested {
		b.invested = true
		b.nextRebalance = b.rebalancePeriod
		goodStocks := getGoodStocks(dailyPrices)
		return gravy.InvestApproximatelyUniformlyInTargets(b.algorithmID, portfolio, dailyPrices, goodStocks)
	} else if b.nextRebalance == 0 {
		b.invested = false
		return gravy.SellEverythingMarketOrder(b.algorithmID, portfolio)
	}
	return nil
}

// New creates a new, uninitialized BuyGoodAndHold algorithm.
func New(algorithmID string, rebalancePeriod int) *BuyGoodAndHold {
	return &BuyGoodAndHold{
		id:              algorithmID,
		algorithmID:     &supervisor_pb.AlgorithmId{AlgorithmId: algorithmID},
		invested:        false,
		rebalancePeriod: rebalancePeriod,
	}
}

// ******************************
//  Mostly boilerplate hereafter
// ******************************

// Init initializes the registrar. The algorithm should be listening before calling Init to avoid deadlocks.
func (b *BuyGoodAndHold) Init() error {
	var err error
	b.registrar, err = registrar.NewWithSupervisor()
	return err
}

// Close closes the regitsrar.
func (b *BuyGoodAndHold) Close() {
	b.registrar.Close()
}

// Execute implements the algorithm interface.
func (b *BuyGoodAndHold) Execute(ctx context.Context, input *algorithmio_pb.Input) (*algorithmio_pb.Output, error) {
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
