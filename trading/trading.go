package trading

import (
	"log"
	"time"

	"github.com/Clever/go-utils/stringset"
)

// Prices is a mapping from asset identifier to a price.
type Prices map[string]float64

// Window represents a unit of time, with various prices for the assets.
type Window struct {
	Begin time.Time
	End   time.Time

	Close  Prices
	High   Prices
	Low    Prices
	Open   Prices
	Volume Prices

	Symbols stringset.StringSet
}

// Portfolio.
type Portfolio struct {
	Stocks  map[string]int
	CashUSD float64
}

func NewPortfolio(seed float64) Portfolio {
	return Portfolio{
		Stocks:  map[string]int{},
		CashUSD: seed,
	}

}

// Returns mature value of a position given a Tick.
func (p Portfolio) Value(prices Prices) float64 {
	value := p.CashUSD
	for symbol, quantity := range p.Stocks {
		value += float64(quantity) * prices[symbol]
	}
	return value
}

// Abstract portfolio. A probability distribution over potential assets.
type AbstractPortfolio struct {
	stocks map[string]float64
	total  float64
}

func NewAbstractPortfolio() AbstractPortfolio {
	return AbstractPortfolio{stocks: map[string]float64{}, total: 0.0}
}

func NewBalancedAbstractPortfolio(stocksPrices Prices) {
	portfolio := NewAbstractPortfolio()
	for symbol, _ := range prices {
		portfolio.SetStock(symbol, 1.0)
	}
	return portfolio
}

// Ensures that the abstract portfolio is normalized to be a distribution.
func (a *AbstractPortfolio) ensureDistribution() {
	if a.total == 0.0 {
		log.Fatalf("Tried to create distribution in AbstractPortfolio before SetStock.")
	} else if a.total == 1.0 {
		return
	}
	for symbol, value := range a.stocks {
		a.stocks[symbol] = value / a.total
	}
	a.total = 1.0
}

// Sets a value for stock and updates the total.
func (a *AbstractPortfolio) SetStock(symbol string, value float64) {
	a.total -= a.stocks[symbol]
	a.stocks[symbol] = value
	a.total += value
}

// Gets stock (after ensuring the portfolio is a distribution).
func (a *AbstractPortfolio) GetStock(symbol string) float64 {
	a.ensureDistribution()
	return a.stocks[symbol]
}

// Balances an abstract portfolio equally into all stocks. Returns any current value for unlisted
// stocks.
func (p *AbstractPortfolio) Balance(prices, previousPrices Prices, seed float64) {
	cashUSD := seed
	for symbol, units := range p.Stocks {
		if previousPrices != nil {
			// For symbols that are no longer listed, divest at previous closing price.
			if _, ok := prices[symbol]; !ok {
				cashUSD += p.Stocks[symbol] * previousPrices[symbol]
			}
		}
		delete(p.Stocks, symbol)
	}

	// Invest equally in all stocks.
	frac := 1.0 / len(prices)
	for symbol, price := range prices {
		b.Stocks[symbol] = frac / price
	}

	return cashout
}

// TODO: include enough data for other types of orders and shorts.
type Order struct {
	StopPrice  float64
	LimitPrice float64
	Volume     int
}

type OrderOutcome struct {
}
