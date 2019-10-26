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

// TODO(desa): not sure what was intended here
//func NewBuyAndHold() *BuyAndHold {
//	return &BuyAndHold{desire: trading.NewAbstractPortfolio(1.0), initialized: false}
//}
