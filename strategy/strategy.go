package strategy

import (
	"time"

	"github.com/d-sparks/gravy/db"
	"github.com/d-sparks/gravy/signal"
	"github.com/d-sparks/gravy/trading"
)

// Strategy outputs about desired position. Will eventually include other data like confidence.
type StrategyOutput struct {
	DesiredPortfolio trading.AbstractPortfolio
}

// Strategies are used in TradingAlgorithm. They represent abstract strategies but don't have the
// ability to figure out which trades to make.
type Strategy interface {
	// Run the strategy.
	Run(
		date time.Time,
		data map[string]db.Store,
		signals map[string]signal.Signal,
	) StrategyOutput

	// Get debug output for previous Run.
	DebugHeaders() []string
	Debug() map[string]string
}
