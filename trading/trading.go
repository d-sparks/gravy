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

func (w Window) MeanHighLowPrice(symbol string) float64 {
	return (w.High[symbol] + w.Low[symbol]) / 2.0
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
func (p *Portfolio) Value(prices Prices) float64 {
	value := p.CashUSD
	for symbol, units := range p.Stocks {
		value += float64(units) * prices[symbol]
	}
	return value
}

// Mean value during a window.
func (p *Portfolio) MeanHighLowValue(window Window) float64 {
	value := p.CashUSD
	c := 0.0
	for symbol, units := range p.Stocks {
		// value += window.Close[symbol] * units //, with Kahan summation
		v := window.MeanHighLowPrice(symbol)*float64(units) - c
		t := value + v
		c = (t - value) - v
		value = t
	}
	return value
}

// AbstractPortfolio. Can have fractional/negative shares.
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
func (a *AbstractPortfolio) Value(prices Prices) float64 {
	value := a.CashUSD
	for symbol, units := range a.Stocks {
		value += float64(units) * prices[symbol]
	}
	return value
}

// Mean value during a window.
func (a *AbstractPortfolio) MeanHighLowValue(window Window) float64 {
	value := a.CashUSD
	c := 0.0
	for symbol, units := range a.Stocks {
		// value += window.Close[symbol] * units //, with Kahan summation
		v := window.MeanHighLowPrice(symbol)*units - c
		t := value + v
		c = (t - value) - v
		value = t
	}
	return value
}

// CapitalDistribution portfolio. A probability distribution over potential assets.
type CapitalDistribution struct {
	stocks map[string]float64
	total  float64
}

func NewCapitalDistribution() CapitalDistribution {
	return CapitalDistribution{stocks: map[string]float64{}, total: 0.0}
}

func NewBalancedCapitalDistribution(prices Prices) CapitalDistribution {
	distribution := NewCapitalDistribution()
	for symbol, _ := range prices {
		distribution.SetStock(symbol, 1.0)
	}
	return distribution
}

// Ensures that the abstract portfolio is normalized to be a distribution.
func (a *CapitalDistribution) ensureDistribution() {
	if a.total == 0.0 {
		log.Fatalf("Tried to create distribution in CapitalDistribution before SetStock.")
	} else if a.total == 1.0 {
		return
	}
	for symbol, value := range a.stocks {
		a.stocks[symbol] = value / a.total
	}
	a.total = 1.0
}

// Sets a value for stock and updates the total.
func (a *CapitalDistribution) SetStock(symbol string, value float64) {
	a.total -= a.stocks[symbol]
	a.stocks[symbol] = value
	a.total += value
}

// Gets stock (after ensuring the portfolio is a distribution).
func (a *CapitalDistribution) GetStock(symbol string) float64 {
	a.ensureDistribution()
	return a.stocks[symbol]
}

// Converts to the equivalent abstract portfolio at given prices.
func (a *CapitalDistribution) ToAbstractPortfolioOnPrices(
	prices Prices,
	value float64,
) AbstractPortfolio {
	a.ensureDistribution()
	portfolio := NewAbstractPortfolio(0.0)
	for symbol, allocation := range a.stocks {
		portfolio.Stocks[symbol] = allocation * value / prices[symbol]
	}
	return portfolio
}

// Expected relative change from window open to close.
func (a *CapitalDistribution) RelativeWindowPerformance(window Window) float64 {
	a.ensureDistribution()
	perf := 0.0
	for symbol, exposure := range a.stocks {
		if exposure == 0.0 {
			continue
		} else if !window.Symbols.Contains(symbol) {
			log.Fatalf("Capital in unlisted symbol")
		}
		open := window.Open[symbol]
		perf += exposure * (window.Close[symbol] - open) / open
	}
	return perf
}

// TODO: include enough data for other types of orders and shorts.
type Order struct {
	StopPrice  float64
	LimitPrice float64
	Volume     int
}

type OrderOutcome struct {
}
