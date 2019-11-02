package buyandhold

import (
	"fmt"
	"strconv"
	"time"

	"github.com/d-sparks/gravy/db"
	"github.com/d-sparks/gravy/db/dailywindow"
	"github.com/d-sparks/gravy/signal"
	"github.com/d-sparks/gravy/signal/unlistings"
	"github.com/d-sparks/gravy/strategy"
	"github.com/d-sparks/gravy/trading"
)

func Name(period int) string {
	return fmt.Sprintf("buyandhold%ddays", period)
}

// BuyAndHold strategy. Invest equally in all securities. Rebalance every N days.
type BuyAndHold struct {
	distribution trading.CapitalDistribution
	debug        map[string]string
	period       int
	days         int

	// Cyclic buffer
	perf       []float64
	periodPerf float64
	wins       int
	losses     int
}

func New(period int) *BuyAndHold {
	return &BuyAndHold{
		period:     period,
		days:       0,
		perf:       make([]float64, period),
		periodPerf: 1.0,
	}
}

// If this is the first time this strategy has run, invest equally in all securities.
func (b *BuyAndHold) Run(
	date time.Time,
	stores map[string]db.Store,
	signals map[string]signal.Signal,
) strategy.StrategyOutput {
	b.debug = map[string]string{}

	// Every Nth day, rebalance to full diversification.
	window := stores[dailywindow.Name].Get(date).Window
	if b.days%b.period == 0 {
		b.distribution = trading.NewBalancedCapitalDistribution(window.Open)
	}

	// Remove allocation to unlisted stocks.
	allUnlistings := signals[unlistings.Name].Compute(date, stores).StringSet
	for symbol, _ := range allUnlistings {
		b.distribution.SetStock(symbol, 0.0)
	}

	// Track performance.
	perf := b.distribution.RelativeWindowPerformance(window)
	if b.days >= b.period {
		b.periodPerf /= (1.0 + b.perf[b.days%b.period])
	}
	b.periodPerf *= (1.0 + perf)
	b.perf[b.days%b.period] = perf
	if perf > 0.0 {
		b.wins++
	} else if perf < 0.0 {
		b.losses++
	}

	// Debug.
	b.debug["periodperf"] = strconv.FormatFloat(b.periodPerf, 'G', -1, 64)
	b.debug["perf"] = strconv.FormatFloat(perf, 'G', -1, 64)
	b.debug["W/L"] = strconv.FormatFloat(float64(b.wins)/float64(b.losses), 'G', -1, 64)

	b.days++
	return strategy.StrategyOutput{CapitalDistribution: &b.distribution}
}

// No data to output
func (b *BuyAndHold) Headers() []string {
	return []string{"perf", "periodperf", "W/L"}
}

// Return last run's debug.
func (b *BuyAndHold) Debug() map[string]string {
	return b.debug
}
