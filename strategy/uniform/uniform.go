package uniform

import (
	"fmt"
	"time"

	"github.com/d-sparks/gravy/db"
	"github.com/d-sparks/gravy/db/dailyprices"
	"github.com/d-sparks/gravy/signal"
	"github.com/d-sparks/gravy/signal/unlistings"
	"github.com/d-sparks/gravy/strategy"
	"github.com/d-sparks/gravy/trading"
)

// Name of the uniform strategy.
const Name = "uniform"

// Uniform strategy. Invest equally in all securities.
type Uniform struct {
	debug map[string]string
}

// New Uniform strategy.
func New() *Uniform {
	return &Uniform{}
}

// Name as per the strategy interface.
func (b *Uniform) Name() string {
	return Name
}

// Run the strategy. Always returns a capital distribution on (today's stocks) \bigcap (tomorrow's stocks) that is
// invested equally in all tickers as today's closing prices.
func (b *Uniform) Run(
	date time.Time,
	stores map[string]db.Store,
	signals map[string]signal.Signal,
) (*strategy.StrategyOutput, error) {
	b.debug = map[string]string{}

	// Get today's prices.
	dailyPricesData, err := stores[dailyprices.Name].Get(date)
	if err != nil {
		return nil, fmt.Errorf("Error getting daily prices in strategy `%s`: `%s`", Name, err.Error())
	}

	// Get unlistings.
	unlistingsData, err := signals[unlistings.Name].Compute(date, stores)
	if err != nil {
		return nil, fmt.Errorf("Error computing unlistings in strategy `%s`: `%s`", Name, err.Error())
	}

	// Create uniform capital distribution.
	investIn := dailyPricesData.Tickers
	for ticker := range unlistingsData.StringSet {
		investIn.Remove(ticker)
	}
	uniform := trading.NewUniformCapitalDistribution(investIn)
	output := strategy.StrategyOutput{CapitalDistribution: uniform}

	return &output, nil
}

// Headers as per strategy interface.
func (b *Uniform) Headers() []string {
	return []string{}
}

// Debug as per strategy interface.
func (b *Uniform) Debug() map[string]string {
	return b.debug
}
