package buythedip

import "github.com/d-sparks/gravy/trading"

const Name = "buythedip"

// BuyTheDip strategy. Attempts to buy underpriced and sell overpriced stocks.
type BuyTheDip struct {
	desire         trading.AbstractPortfolio
	debug          map[string]string
	previousWindow trading.Window
	initialized    bool
}

func NewBuyAndHold() *BuyAndHold {
	return &BuyAndHold{desire: trading.NewAbstractPortfolio(1.0), initialized: false}
}
