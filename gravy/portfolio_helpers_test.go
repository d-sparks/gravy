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

var testAlgorithmID = &supervisor_pb.AlgorithmId{AlgorithmId: "TEST"}

// TestInvestApproximatelyUniformlyHappyPath is the happy path test case with a single stock.
func TestInvestApproximatelyUniformlyHappyPath(t *testing.T) {
	var dailyPrices dailyprices_pb.DailyPrices
	var portfolio supervisor_pb.Portfolio

	dailyPrices.StockPrices = map[string]*dailyprices_pb.DailyPrices_StockPrices{
		"MSFT": &dailyprices_pb.DailyPrices_StockPrices{Close: 10.0},
	}
	portfolio.Usd = 1000.0

	orders := InvestApproximatelyUniformly(testAlgorithmID, &portfolio, &dailyPrices)

	assert.Equal(t, 1, len(orders))
	order := orders[0]
	assert.Equal(t, "MSFT", order.GetTicker())
	assert.Equal(t, 99.0, order.GetVolume())
	assert.Equal(t, 10.1, order.GetLimit())
}

// TestInvestApproximatelyUniformlyTwoStocks is the happy path test case with a two stocks and no tradeoff.
func TestInvestApproximatelyUniformlyTwoStocks(t *testing.T) {
	var dailyPrices dailyprices_pb.DailyPrices
	var portfolio supervisor_pb.Portfolio

	dailyPrices.StockPrices = map[string]*dailyprices_pb.DailyPrices_StockPrices{
		"MSFT": &dailyprices_pb.DailyPrices_StockPrices{Close: 10.0},
		"GOOG": &dailyprices_pb.DailyPrices_StockPrices{Close: 20.0},
	}
	portfolio.Usd = 1000.0

	orders := InvestApproximatelyUniformly(testAlgorithmID, &portfolio, &dailyPrices)

	expectedVolume := map[string]float64{"MSFT": 49.0, "GOOG": 25.0}
	expectedLimit := map[string]float64{"MSFT": 10.1, "GOOG": 20.2}
	assert.Equal(t, 3, len(orders))
	volume := map[string]float64{}
	for _, order := range orders {
		assert.Equal(t, expectedLimit[order.GetTicker()], order.GetLimit())
		volume[order.GetTicker()] += order.GetVolume()
	}

	assert.Equal(t, expectedVolume, volume)
}

// TestInvestApproximatelyUniformlyManyStocks tests the case of many stocks.
func TestInvestApproximatelyUniformlyManyStocks(t *testing.T) {
	var dailyPrices dailyprices_pb.DailyPrices
	var portfolio supervisor_pb.Portfolio

	dailyPrices.StockPrices = map[string]*dailyprices_pb.DailyPrices_StockPrices{
		"MSFT": &dailyprices_pb.DailyPrices_StockPrices{Close: 10.0},
		"GOOG": &dailyprices_pb.DailyPrices_StockPrices{Close: 20.0},
		"FB":   &dailyprices_pb.DailyPrices_StockPrices{Close: 7.0},
		"APPL": &dailyprices_pb.DailyPrices_StockPrices{Close: 150.0},
		"NVDA": &dailyprices_pb.DailyPrices_StockPrices{Close: 9.0},
		"GM":   &dailyprices_pb.DailyPrices_StockPrices{Close: 0.50},
		"FORD": &dailyprices_pb.DailyPrices_StockPrices{Close: 1.0},
	}
	portfolio.Usd = 1000.0

	orders := InvestApproximatelyUniformly(testAlgorithmID, &portfolio, &dailyPrices)

	expectedVolume := map[string]float64{
		"APPL": 1.0,
		"FB":   20.0,
		"FORD": 141.0,
		"GM":   282.0,
		"GOOG": 7.0,
		"MSFT": 14.0,
		"NVDA": 15.0,
	}
	expectedLimit := map[string]float64{
		"APPL": 151.5,
		"FB":   7.07,
		"FORD": 1.01,
		"GM":   0.505,
		"GOOG": 20.2,
		"MSFT": 10.1,
		"NVDA": 9.09,
	}

	assert.Equal(t, 7, len(orders))
	volume := map[string]float64{}
	limit := map[string]float64{}
	totalLimit := 0.0
	for _, order := range orders {
		volume[order.GetTicker()] += order.GetVolume()
		limit[order.GetTicker()] = order.GetLimit()
		totalLimit += order.GetLimit() * order.GetVolume()
	}

	assert.Equal(t, expectedVolume, volume)
	assert.Equal(t, expectedLimit, limit)
	assert.Less(t, 996.0, totalLimit)
}
