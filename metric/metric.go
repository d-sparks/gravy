package metric

import (
	"time"

	"github.com/d-sparks/gravy/db"
	"github.com/d-sparks/gravy/signal"
	"github.com/d-sparks/gravy/strategy"
)

// This is possibly temporary until we figure out what metrics are in general.
type PerStrategyMetric interface {
	// Identifier of metric.
	Name() string

	// Value of metric on a given date.
	Value(
		date time.Time,
		stores map[string]db.Store,
		signals map[string]signal.Signal,
		strategyOutput *strategy.StrategyOutput,
	) (float64, error)
}
