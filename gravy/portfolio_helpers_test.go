package gravy

import (
	"testing"

	dailyprices_pb "github.com/d-sparks/gravy/data/dailyprices/proto"
	supervisor_pb "github.com/d-sparks/gravy/supervisor/proto"
	"github.com/stretchr/testify/assert"
)

var (
	algorithmID = &supervisor_pb.AlgorithmId{AlgorithmId: "test"}

	portfolio = &supervisor_pb.Portfolio{Stocks: map[string]float64{"MSFT": 3.0, "GOOG": 2.0}, Usd: 1.0}

	prices = &dailyprices_pb.DailyPrices{
		StockPrices: map[string]*dailyprices_pb.DailyPrices_StockPrices{
			"MSFT": &dailyprices_pb.DailyPrices_StockPrices{Close: 5.0},
			"GOOG": &dailyprices_pb.DailyPrices_StockPrices{Close: 10.0},
			"FB":   &dailyprices_pb.DailyPrices_StockPrices{Close: 7.0},
		},
	}
)

// TestPortfolioValue tests the basic arithmetic of portfolio value.
func TestPortfolioValue(t *testing.T) {
	assert.Equal(t, 36.0, PortfolioValue(portfolio, prices))
}

// TestTargetUniformInvestment tests the arithmetic. We expect portfolio value / # stocks.
func TestTargetUniformInvestment(t *testing.T) {
	assert.Equal(t, 12.0, TargetUniformInvestment(portfolio, prices))
}

// TestSellEverythingWithStop tests with custom stop lambda.
func TestSellEverythingWithStop(t *testing.T) {
	stop := func(ticker string) float64 {
		return map[string]float64{"MSFT": 0.1, "GOOG": 0.2}[ticker]
	}
	orders := SellEverythingWithStop(algorithmID, portfolio, stop)

	expected := map[string]*supervisor_pb.Order{
		"MSFT": &supervisor_pb.Order{
			AlgorithmId: algorithmID, Ticker: "MSFT", Volume: -3.0, Stop: 0.1, Limit: 0.0,
		},
		"GOOG": &supervisor_pb.Order{
			AlgorithmId: algorithmID, Ticker: "GOOG", Volume: -2.0, Stop: 0.2, Limit: 0.0,
		},
	}
	for _, order := range orders {
		assert.Equal(t, expected[order.GetTicker()], order)
	}
}

// TestSellEverythingWithStopPercent tests with fixed stop percent.
func TestSellEverythingWithStopPercent(t *testing.T) {
	orders := SellEverythingWithStopPercent(algorithmID, portfolio, prices, 0.5)

	expected := map[string]*supervisor_pb.Order{
		"MSFT": &supervisor_pb.Order{
			AlgorithmId: algorithmID, Ticker: "MSFT", Volume: -3.0, Stop: 2.5, Limit: 0.0,
		},
		"GOOG": &supervisor_pb.Order{
			AlgorithmId: algorithmID, Ticker: "GOOG", Volume: -2.0, Stop: 5.0, Limit: 0.0,
		},
	}
	for _, order := range orders {
		assert.Equal(t, expected[order.GetTicker()], order)
	}
}

// TestSellEverythingWithMarketOrder tests with no stop.
func TestSellEverythingMarketOrder(t *testing.T) {
	orders := SellEverythingMarketOrder(algorithmID, portfolio)

	expected := map[string]*supervisor_pb.Order{
		"MSFT": &supervisor_pb.Order{
			AlgorithmId: algorithmID, Ticker: "MSFT", Volume: -3.0, Stop: 0.0, Limit: 0.0,
		},
		"GOOG": &supervisor_pb.Order{
			AlgorithmId: algorithmID, Ticker: "GOOG", Volume: -2.0, Stop: 0.0, Limit: 0.0,
		},
	}
	for _, order := range orders {
		assert.Equal(t, expected[order.GetTicker()], order)
	}
}
