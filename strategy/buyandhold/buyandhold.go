package buyandhold

import (
	"fmt"
	"time"

	"github.com/d-sparks/gravy/db"
	"github.com/d-sparks/gravy/db/dailywindow"
	"github.com/d-sparks/gravy/signal"
	"github.com/d-sparks/gravy/strategy"
	"github.com/d-sparks/gravy/trading"
)

const Name = "buyandhold"

// BuyAndHold strategy. Invest equally in all securities and hold them forever.
type BuyAndHold struct {
	desire         trading.AbstractPortfolio
	debug          map[string]string
	previousWindow trading.Window
	initialized    bool
}

func NewBuyAndHold() *BuyAndHold {
	return &BuyAndHold{desire: trading.NewAbstractPortfolio(1.0), initialized: false}
}

// If this is the first time this strategy has run, invest equally in all securities.
func (b *BuyAndHold) Run(
	date time.Time,
	stores map[string]db.Store,
	signals map[string]signal.Signal,
) strategy.StrategyOutput {
	b.debug = map[string]string{}
	window := stores[dailywindow.Name].Get(date).Window

	if !b.initialized {
		// Initialize by investing seed equally in all available stocks.
		frac := b.desire.CashUSD / float64(len(window.Open))
		for symbol, _ := range window.Open {
			b.desire.Stocks[symbol] = frac / window.Close[symbol]
		}
		b.desire.CashUSD = 0.0
		b.initialized = true
	} else {
		// For symbols that are no longer listed, divest at previous closing price.
		for symbol, closePrice := range b.previousWindow.Close {
			if _, ok := window.Open[symbol]; !ok {
				b.desire.CashUSD += b.desire.Stocks[symbol] * closePrice
				delete(b.desire.Stocks, symbol)
			}
		}
	}

	b.previousWindow = window
	b.debug["value"] = fmt.Sprintf("%f", b.desire.Value(window.Close))

	return strategy.StrategyOutput{DesiredPortfolio: b.desire}
}

// No data to output
func (b *BuyAndHold) Headers() []string {
	return []string{"value"}
}

// Return last run's debug.
func (b *BuyAndHold) Debug() map[string]string {
	return b.debug
}
