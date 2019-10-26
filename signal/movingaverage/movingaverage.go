package movingaverage

import (
	"fmt"
	"time"

	"github.com/d-sparks/gravy/db"
	"github.com/d-sparks/gravy/db/dailywindow"
	"github.com/d-sparks/gravy/signal"
)

func Name(days int) string { return fmt.Sprintf("%dday_movingaverage", days) }

// Tracks moving average and standard deviation over N days.
type MovingAverage struct {
	observations map[string][]float64
	oldestIx     map[string]int
	days         int
}

func NewMovingAverage(days int) *MovingAverage {
	return signal.NewCachedSignal(
		&MovingAverage{
			days:         days,
			oldestIx:     map[string]int{},
			observations: map[string][]float64{},
		},
		time.Hour*24,
	)
}

func (m *MovingAverage) Compute(date time.Time, stores map[string]db.Store) signal.SignalOutput {
	// Get newest window
	window := stores[dailywindow.Name].Get(date).Window

	// Stop tracking unlisted symbols.
	for symbol, _ := range m.observations {
		if _, ok := window.Close[symbol]; !ok {
			delete(m.observations, symbol)
			delete(m.oldestIx, symbol)
		}
	}

	// Record new observations.
	for symbol, price := range window.Close {
		if len(m.observations[symbol]) < m.days {
			m.observations[symbol] = append(m.observations[symbol], price)
			continue
		}
		m.observations[symbol][m.oldestIx[symbol]] = price
		m.oldestIx[symbol] = (m.oldestIx[symbol] + 1) % m.days
	}

	// Update observations...
	return signal.SignalOutput{}
}

func (m *MovingAverage) Headers() []string {
	return []string{}
}

func (m *MovingAverage) Debug() map[string]string {
	return map[string]string{}
}
