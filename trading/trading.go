package trading

import (
	"time"
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

// Abstract portfolio. Can hold fractional and negative shares.
type AbstractPortfolio struct {
	Stocks  map[string]float64
	CashUSD float64
}

func NewAbstractPortfolio(seed float64) AbstractPortfolio {
	return AbstractPortfolio{
		Stocks:  map[string]float64{},
		CashUSD: seed,
	}

}

// Returns mature value of a position given a Tick.
func (p AbstractPortfolio) Value(prices Prices) float64 {
	value := p.CashUSD
	for symbol, quantity := range p.Stocks {
		value += quantity * prices[symbol]
	}
	return value
}

// TODO: include enough data for other types of orders and shorts.
type Order struct {
	StopPrice  float64
	LimitPrice float64
	Volume     int
}

type OrderOutcome struct {
}
