package ipos

import (
	"time"

	"github.com/Clever/go-utils/stringset"
	"github.com/d-sparks/gravy/db"
	"github.com/d-sparks/gravy/db/dailywindow"
	"github.com/d-sparks/gravy/signal"
)

const Name = "ipos"

// IPOs reports new values since last tick.
type IPOs struct {
	previousSymbols stringset.StringSet
}

func New() IPOs {
	return signal.NewCachedSignal(
		IPOs{previousSymbols: stringset.StringSet{}},
		time.Hour*24*365, // one year
	)
}

func (i *IPOs) Compute(date time.Time, stores map[string]db.Store) signal.SignalOutput {
	symbols := stores[dailywindow.Name].Get(date).Window.Symbols
	output := signal.SignalOutput{StringSet: symbols.Minus(i.previousSymbols)}
	i.previousSymbols = symbols
	return output
}

func (i *IPOs) Headers() []string {
	return []string{}
}

func (i *IPOs) Debug() map[string]string {
	return map[string]string{}
}
