package ipos

import (
	"fmt"
	"time"

	"github.com/Clever/go-utils/stringset"
	"github.com/d-sparks/gravy/db"
	"github.com/d-sparks/gravy/db/dailyprices"
	"github.com/d-sparks/gravy/signal"
)

const Name = "ipos"

// IPOs reports tickers not present in the last daily prices.
type IPOs struct {
	previousTickers stringset.StringSet
}

func NewIPOs() *signal.CachedSignal {
	return signal.NewCachedSignal(
		&IPOs{previousTickers: stringset.StringSet{}},
		time.Hour*24*365, // one year
	)
}

func (i *IPOs) Name() string {
	return Name
}

// Diffs tickers reported by dailyprices with that seen in the previous computation.
func (i *IPOs) Compute(date time.Time, stores map[string]db.Store) (*signal.SignalOutput, error) {
	dailyPrices, err := stores[dailyprices.Name].Get(date)
	if err != nil {
		return nil, fmt.Errorf("Error reading dailyprices in `%s`: `%s`", Name, err.Error())
	}
	output := signal.SignalOutput{StringSet: dailyPrices.Tickers.Minus(i.previousTickers)}
	i.previousTickers = dailyPrices.Tickers
	return &output, nil
}

func (i *IPOs) Headers() []string {
	return []string{}
}

func (i *IPOs) Debug() map[string]string {
	return map[string]string{}
}
