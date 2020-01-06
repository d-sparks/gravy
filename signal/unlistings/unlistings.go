package unlistings

import (
	"fmt"
	"time"

	"github.com/Clever/go-utils/stringset"
	"github.com/d-sparks/gravy/db"
	"github.com/d-sparks/gravy/db/dailyprices"
	"github.com/d-sparks/gravy/signal"
)

const Name = "unlistings"

// Unlistings reports new values since last tick.
type Unlistings struct {
	previousTickers stringset.StringSet
}

func NewUnlistings() *signal.CachedSignal {
	return signal.NewCachedSignal(
		&Unlistings{previousTickers: stringset.StringSet{}},
		time.Hour*24*365,
	)
}

func (u *Unlistings) Name() string {
	return Name
}

func (u *Unlistings) Compute(date time.Time, stores map[string]db.Store) (*signal.SignalOutput, error) {
	data, err := stores[dailyprices.Name].Get(date)
	if err != nil {
		return nil, fmt.Errorf("Error reading dailyprices in `%s`: `%s`", Name, err.Error())
	}
	output := signal.SignalOutput{StringSet: u.previousTickers.Minus(data.Tickers)}
	u.previousTickers = data.Tickers
	return &output, nil
}

func (u *Unlistings) Headers() []string {
	return []string{}
}

func (u *Unlistings) Debug() map[string]string {
	return map[string]string{}
}
