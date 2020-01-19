package unlistings

import (
	"fmt"
	"time"

	"github.com/d-sparks/gravy/db"
	"github.com/d-sparks/gravy/db/dailyprices"
	"github.com/d-sparks/gravy/signal"
)

// Name of the unlistings signal.
const Name = "unlistings"

// Unlistings reports tickers that will not be listed tomorrow.
type Unlistings struct{}

// New cached Unlistings.
func New() *signal.CachedSignal {
	return signal.NewCachedSignal(&Unlistings{}, time.Hour*24*365)
}

// Name returns the name of Unlistings as per the signal interface.
func (u *Unlistings) Name() string {
	return Name
}

// Compute computes the signal as per the signal interface.
func (u *Unlistings) Compute(date time.Time, stores map[string]db.Store) (*signal.SignalOutput, error) {
	todayData, err := stores[dailyprices.Name].Get(date)
	if err != nil {
		return nil, fmt.Errorf("Error reading todayprices in `%s`: `%s`", Name, err.Error())
	}

	tomorrowDate, err := stores[dailyprices.Name].NextDate(date)
	if err != nil {
		return nil, fmt.Errorf("Error getting next date in `%s`: `%s`", Name, err.Error())
	}

	tomorrowData, err := stores[dailyprices.Name].Get(*tomorrowDate)
	if err != nil {
		return nil, fmt.Errorf("Error getting tomorrowprices in `%s`: `%s`", Name, err.Error())
	}
	return &signal.SignalOutput{StringSet: todayData.Tickers.Minus(tomorrowData.Tickers)}, nil
}

// Headers as per signal interface.
func (u *Unlistings) Headers() []string {
	return []string{}
}

// Debug as per signal interface.
func (u *Unlistings) Debug() map[string]string {
	return map[string]string{}
}
