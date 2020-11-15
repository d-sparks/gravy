package buyandhold

import (
	"testing"

	dailyprices_pb "github.com/d-sparks/gravy/data/dailyprices/proto"
	supervisor_pb "github.com/d-sparks/gravy/supervisor/proto"
	"github.com/stretchr/testify/assert"
)

// TestInvestApproximatelyUniformlyHappyPath is the happy path test case with a single stock.
func TestInvestApproximatelyUniformlyHappyPath(t *testing.T) {
	var dailyPrices dailyprices_pb.DailyPrices
	var portfolio supervisor_pb.Portfolio

	dailyPrices.StockPrices = map[string]*dailyprices_pb.DailyPrices_StockPrices{
		"MSFT": &dailyprices_pb.DailyPrices_StockPrices{Close: 10.0},
	}
	portfolio.Usd = 1000.0

	b := New()
	orders := b.InvestApproximatelyUniformly(&portfolio, &dailyPrices)

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

	b := New()
	orders := b.InvestApproximatelyUniformly(&portfolio, &dailyPrices)

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

	b := New()
	orders := b.InvestApproximatelyUniformly(&portfolio, &dailyPrices)

	expectedVolume := map[string]float64{
		"APPL": 1.0,
		"FB":   20.0,
		"FORD": 144.0,
		"GM":   282.0,
		"GOOG": 7.0,
		"MSFT": 14.0,
		"NVDA": 15,
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

	assert.Equal(t, 10, len(orders))
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
	assert.Less(t, 999.0, totalLimit)
}
