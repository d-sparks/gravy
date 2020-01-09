package movingaverage

import (
	"fmt"
	"time"

	"github.com/d-sparks/gravy/db"
	"github.com/d-sparks/gravy/db/dailyprices"
	"github.com/d-sparks/gravy/signal"
	"github.com/montanaflynn/stats"
)

func Name(days int) string { return fmt.Sprintf("%dday_movingaverage", days) }

// Tracks moving average and standard deviation over N days.
type MovingAverage struct {
	observations map[string][]float64
	oldestIx     map[string]int
	days         int
}

func New(days int) *signal.CachedSignal {
	return signal.NewCachedSignal(
		&MovingAverage{
			days:         days,
			oldestIx:     map[string]int{},
			observations: map[string][]float64{},
		},
		time.Hour*24,
	)
}

func (m *MovingAverage) Name() string {
	return Name(m.days)
}

func (m *MovingAverage) Compute(date time.Time, stores map[string]db.Store) (*signal.SignalOutput, error) {
	data, err := stores[dailyprices.Name].Get(date)
	if err != nil {
		return nil, fmt.Errorf("Error reading dailyprices in `%s`: `%s`", Name(m.days), err.Error())
	}

	// Stop tracking unlisted tickers.
	for ticker, _ := range m.observations {
		if !data.Tickers.Contains(ticker) {
			delete(m.observations, ticker)
			delete(m.oldestIx, ticker)
		}
	}

	// Record new observations.
	for ticker, prices := range data.TickersToPrices {
		if len(m.observations[ticker]) < m.days {
			m.observations[ticker] = append(m.observations[ticker], prices.Close)
			continue
		}
		m.observations[ticker][m.oldestIx[ticker]] = prices.Close
		m.oldestIx[ticker] = (m.oldestIx[ticker] + 1) % m.days
	}

	// Update observations.
	output := signal.SignalOutput{KV: map[string]float64{}}
	for ticker, prices := range m.observations {
		output.KV[ticker], _ = stats.Mean(prices)
	}
	return &output, nil
}

func (m *MovingAverage) Headers() []string {
	return []string{}
}

func (m *MovingAverage) Debug() map[string]string {
	return map[string]string{}
}
