package buyandhold

import (
	"testing"
	"time"

	"github.com/Clever/go-utils/stringset"
	"github.com/d-sparks/gravy/db"
	"github.com/d-sparks/gravy/db/dailyprices"
	"github.com/d-sparks/gravy/signal"
	"github.com/d-sparks/gravy/signal/ipos"
	"github.com/d-sparks/gravy/signal/unlistings"
	"github.com/stretchr/testify/assert"
)

func TestBuyAndHoldExample(t *testing.T) {
	// Load some simple fake data.
	store := dailyprices.NewInMemoryStore()
	parsed := map[string]time.Time{}
	pricesMap := map[string]map[string]float64{
		"2020-01-01": map[string]float64{"MSFT": 1.0},
		"2020-01-02": map[string]float64{"MSFT": 2.0},
		"2020-01-03": map[string]float64{"MSFT": 1.0, "GOOGL": 1.0},
		"2020-01-04": map[string]float64{"MSFT": 1.0, "GOOGL": 2.0},
		"2020-01-05": map[string]float64{"MSFT": 2.0, "GOOGL": 2.0},
		"2020-01-06": map[string]float64{"MSFT": 2.0, "GOOGL": 2.0, "FB": 2.0},
		"2020-01-07": map[string]float64{"MSFT": 2.0, "GOOGL": 2.0, "FB": 4.0},
		"2020-01-08": map[string]float64{"MSFT": 2.0, "GOOGL": 2.0, "FB": 4.0},
	}
	for dateStr, prices := range pricesMap {
		date, err := time.Parse("2006-01-02", dateStr)
		assert.NoError(t, err)
		data := db.Data{Tickers: stringset.New(), TickersToPrices: map[string]db.Prices{}}
		for ticker, price := range prices {
			data.Tickers.Add(ticker)
			data.TickersToPrices[ticker] = db.Prices{Close: price}
		}
		store.Set(date, &data)
		parsed[dateStr] = date
	}

	// Create inputs necessary for strategy.
	stores := map[string]db.Store{}
	stores[dailyprices.Name] = store
	signals := map[string]signal.Signal{}
	signals[ipos.Name] = ipos.New()
	signals[unlistings.Name] = unlistings.New()

	// Run strategy on given dates.
	b := New()

	// In the first window, should invest everything in MSFT.
	output, err := b.Run(parsed["2020-01-01"], stores, signals)
	assert.NoError(t, err)
	assert.Equal(t, 1.0, output.CapitalDistribution.GetStock("MSFT"))

	// Second window we are still all-in on MSFT.
	output, err = b.Run(parsed["2020-01-02"], stores, signals)
	assert.NoError(t, err)
	assert.Equal(t, 1.0, output.CapitalDistribution.GetStock("MSFT"))

	// Third window we rebalance to be equal in MSFT and GOOGL.
	output, err = b.Run(parsed["2020-01-03"], stores, signals)
	assert.NoError(t, err)
	assert.Equal(t, 0.5, output.CapitalDistribution.GetStock("MSFT"))
	assert.Equal(t, 0.5, output.CapitalDistribution.GetStock("GOOGL"))

	// Fourth window it's as if we hold one unit of MSFT and one unit of GOOGL.
	output, err = b.Run(parsed["2020-01-04"], stores, signals)
	assert.NoError(t, err)
	assert.Equal(t, 1.0/3.0, output.CapitalDistribution.GetStock("MSFT"))
	assert.Equal(t, 2.0/3.0, output.CapitalDistribution.GetStock("GOOGL"))

	// Fifth window is still one each of MSFT and GOOGL.
	output, err = b.Run(parsed["2020-01-05"], stores, signals)
	assert.NoError(t, err)
	assert.Equal(t, 0.5, output.CapitalDistribution.GetStock("MSFT"))
	assert.Equal(t, 0.5, output.CapitalDistribution.GetStock("GOOGL"))

	// Sixth window is split between MSFT/GOOGL/FB.
	output, err = b.Run(parsed["2020-01-06"], stores, signals)
	assert.NoError(t, err)
	assert.Equal(t, 1.0/3.0, output.CapitalDistribution.GetStock("MSFT"))
	assert.Equal(t, 1.0/3.0, output.CapitalDistribution.GetStock("GOOGL"))
	assert.Equal(t, 1.0/3.0, output.CapitalDistribution.GetStock("FB"))

	// Seventh window is as if the 1/3 of our portfolio doubled (FB). Thus half of our wealth is in FB.
	output, err = b.Run(parsed["2020-01-07"], stores, signals)
	assert.NoError(t, err)
	assert.Equal(t, 0.25, output.CapitalDistribution.GetStock("MSFT"))
	assert.Equal(t, 0.25, output.CapitalDistribution.GetStock("GOOGL"))
	assert.Equal(t, 0.5, output.CapitalDistribution.GetStock("FB"))
}
