package movingaverage

import (
	"testing"
	"time"

	"github.com/Clever/go-utils/stringset"
	"github.com/d-sparks/gravy/db"
	"github.com/d-sparks/gravy/db/dailyprices"
	"github.com/stretchr/testify/assert"
)

func TestMovingAverages(t *testing.T) {
	testData := map[string]map[string]float64{
		"2020-01-06": map[string]float64{"MSFT": 0.0},
		"2020-01-07": map[string]float64{"MSFT": 24.0, "GOOGL": 0.0},
		"2020-01-08": map[string]float64{"MSFT": 48.0, "GOOGL": -24.0},
		"2020-01-09": map[string]float64{"MSFT": 48.0, "GOOGL": -48.0},
		"2020-01-10": map[string]float64{"MSFT": 48.0, "GOOGL": -48.0},
	}
	asDate := map[string]time.Time{}

	testStore := dailyprices.NewInMemoryStore()
	for dateStr, tickersToPrices := range testData {
		date, err := time.Parse("2006-01-02", dateStr)
		assert.NoError(t, err)
		asDate[dateStr] = date
		data := db.Data{TickersToPrices: map[string]db.Prices{}, Tickers: stringset.New()}
		for ticker, price := range tickersToPrices {
			data.TickersToPrices[ticker] = db.Prices{Close: price}
			data.Tickers.Add(ticker)
		}
		testStore.Set(date, &data)
	}

	stores := map[string]db.Store{}
	stores[dailyprices.Name] = testStore

	movingaverage := NewMovingAvareage(3)

	assert.Equal(t, "3day_movingaverage", movingaverage.Name())

	jan6output, err := movingaverage.Compute(asDate["2020-01-06"], stores)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(jan6output.KV))
	assert.Equal(t, 0.0, jan6output.KV["MSFT"])

	jan7output, err := movingaverage.Compute(asDate["2020-01-07"], stores)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(jan7output.KV))
	assert.Equal(t, 12.0, jan7output.KV["MSFT"])
	assert.Equal(t, 0.0, jan7output.KV["GOOGL"])

	jan8output, err := movingaverage.Compute(asDate["2020-01-08"], stores)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(jan8output.KV))
	assert.Equal(t, 24.0, jan8output.KV["MSFT"])
	assert.Equal(t, -12.0, jan8output.KV["GOOGL"])

	jan9output, err := movingaverage.Compute(asDate["2020-01-09"], stores)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(jan9output.KV))
	assert.Equal(t, 40.0, jan9output.KV["MSFT"])
	assert.Equal(t, -24.0, jan9output.KV["GOOGL"])

	jan10output, err := movingaverage.Compute(asDate["2020-01-10"], stores)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(jan10output.KV))
	assert.Equal(t, 48.0, jan10output.KV["MSFT"])
	assert.Equal(t, -40.0, jan10output.KV["GOOGL"])
}
