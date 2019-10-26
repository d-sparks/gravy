package strategy

import (
	"time"

	"github.com/d-sparks/gravy/db"
	"github.com/d-sparks/gravy/signal"
	"github.com/d-sparks/gravy/trading"
)

// Strategy outputs about recommended position. Will eventually include other data like confidence.
type StrategyOutput struct {
	CapitalDistribution *trading.CapitalDistribution
}

// Strategies are used in TradingAlgorithm. They represent abstract strategies but don't have the
// ability to figure out which trades to make.
type Strategy interface {
	// Run the strategy.
	Run(
		date time.Time,
		stores map[string]db.Store,
		signals map[string]signal.Signal,
	) StrategyOutput

	// Get debug output for previous Run.
	Headers() []string
	Debug() map[string]string
}
