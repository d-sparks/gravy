package buy_and_hold

import (
	"github.com/d-sparks/ace-of-trades/trading"
)

// BuyAndHold strategy. Invest equally in all securities.
type BuyAndHold struct {
	Position trading.Position
}

// No data to output
func (b *BuyAndHold) Headers() []string {
	return nil
}

// Invest equally in all securities.
func (b *BuyAndHold) Initialize(tick trading.Tick) (trading.Position, []string) {
	investment := b.Position.Liquid / float64(len(tick))
	for symbol, window := range tick {
		b.Position.Investments[symbol] = investment / window.MeanPrice()
	}
	b.Position.Liquid = 0.0
	return b.Position, nil
}

// Hold, unless a security is no longer listed. In which case, liquidate that security.
func (b *BuyAndHold) ProcessTick(
	tick trading.Tick,
	ipo,
	unlist []string,
	returns float64,
) (trading.Position, []string) {
	for _, symbol := range unlist {
		delete(b.Position.Investments, symbol)
	}
	// Simple strategy: don't reinvest the returns.
	b.Position.Liquid += returns
	return b.Position, nil
}
