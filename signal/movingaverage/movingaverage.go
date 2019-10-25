package signal

import (
	"fmt"
	"time"

	"github.com/d-sparks/gravy/db"
	"github.com/d-sparks/gravy/db/dailywindow"
	"github.com/d-sparks/gravy/trading"
)

func Name(days int) string { return fmt.Sprintf("%dday_movingaverage", days) }

// Tracks moving average and standard deviation over N days.
type MovingAverage struct {
	observations trading.Prices
	oldestIx     map[string]int
	days         int
}

func NewMovingAverage(days int) MovingAverage {
	return MovingAverage{days: days}
}

func (m *MovingAverage) Compute(date time.Time, stores map[string]db.Store) SignalOutput {
	// Get newest window
	window := stores[dailywindow.Name].Get(date).Window

	// Update observations...
}

func (m *MovingAverage) Headers() []string {
	return nil
}

func (m *MovingAverage) Debug() map[string]string {
	return nil
}
