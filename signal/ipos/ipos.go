package ipos

import (
	"fmt"
	"time"

	"github.com/d-sparks/gravy/db"
	"github.com/d-sparks/gravy/db/dailyprices"
	"github.com/d-sparks/gravy/signal"
)

// Name of the ipos signal.
const Name = "ipos"

// IPOs returns symbols which will be be IPOing on the NEXT trading day.
type IPOs struct{}

// New cached IPOs.
func New() *signal.CachedSignal {
	return signal.NewCachedSignal(&IPOs{}, time.Hour*24*365)
}

// Name returns the ipos signal name, for the signal interface.
func (i *IPOs) Name() string {
	return Name
}

// Compute computes the signal as per the signal interface.
func (i *IPOs) Compute(date time.Time, stores map[string]db.Store) (*signal.SignalOutput, error) {
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
	return &signal.SignalOutput{StringSet: tomorrowData.Tickers.Minus(todayData.Tickers)}, nil
}

// Headers as per signal interface.
func (i *IPOs) Headers() []string {
	return []string{}
}

// Debug as per signal interface.
func (i *IPOs) Debug() map[string]string {
	return map[string]string{}
}
