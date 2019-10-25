package signal

import (
	"time"

	"github.com/d-sparks/gravy/db"
)

// Data output by signals.
type SignalOutput struct {
	KV map[string]float64
}

// Signals compute data to be used to inform a Strategy or TradingAlgorithm.
type Signal interface {
	// Compute and/or return cached signal output.
	Compute(date time.Time, stores map[string]db.Store) SignalOutput

	// Debug info for previous computation.
	Headers() []string
	Debug() map[string]string
}
