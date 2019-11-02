package unlistings

import (
	"time"

	"github.com/Clever/go-utils/stringset"
	"github.com/d-sparks/gravy/db"
	"github.com/d-sparks/gravy/db/dailywindow"
	"github.com/d-sparks/gravy/signal"
)

const Name = "unlistings"

// Unlistings reports new values since last tick.
type Unlistings struct {
	previousSymbols stringset.StringSet
}

func New() *signal.CachedSignal {
	return signal.NewCachedSignal(
		&Unlistings{previousSymbols: stringset.StringSet{}},
		time.Hour*24*365,
	)
}

func (u *Unlistings) Compute(date time.Time, stores map[string]db.Store) signal.SignalOutput {
	symbols := stores[dailywindow.Name].Get(date).Window.Symbols
	output := signal.SignalOutput{StringSet: u.previousSymbols.Minus(symbols)}
	u.previousSymbols = symbols
	return output
}

func (u *Unlistings) Headers() []string {
	return []string{}
}

func (u *Unlistings) Debug() map[string]string {
	return map[string]string{}
}
