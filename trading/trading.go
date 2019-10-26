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

func (w Window) MeanPrice(symbol string) float64 {
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
func (p *Portfolio) MeanValue(window Window) float64 {
	value := p.CashUSD
	for symbol, units := range p.Stocks {
		value += window.MeanPrice(symbol) * float64(units)
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
func (a *AbstractPortfolio) MeanValue(window Window) float64 {
	sum := a.CashUSD
	c := 0.0
	for symbol, units := range a.Stocks {
		// sum += window.Close[symbol] * units //, with Kahan summation
		v := window.MeanPrice(symbol)*units - c
		t := sum + v
		c = (t - sum) - v
		sum = t

	}
	return sum
}

// (Re)balance the stock investments of an abstract portfolio.
func (a *AbstractPortfolio) RebalanceOvernight(window, previousWindow Window) float64 {
	previousValue := a.MeanValue(previousWindow)
	// "Sell" any stock currently tracked at the mean window price.
	reinvest := 0.0
	if previousWindow.Symbols == nil {
		reinvest += a.CashUSD
		a.CashUSD = 0.0
	}
	for symbol, shares := range a.Stocks {
		if window.Symbols.Contains(symbol) {
			reinvest += window.Close[symbol] * shares
		} else {
			a.CashUSD += previousWindow.Close[symbol]
		}
		delete(a.Stocks, symbol)
	}
	// "Buy" stock, but not reinvesting for unlisted symbols.
	for symbol, _ := range window.Symbols {
		price := window.Close[symbol]
		a.Stocks[symbol] = reinvest / (price * float64(len(window.Symbols)))
	}

	return a.MeanValue(window) - previousValue
}

// Returns a capital distribution for an abstract portfolio based on mean price over a window.
func (a *AbstractPortfolio) ToCapitalDistributionOnMeanPrice(
	window Window,
) *CapitalDistribution {
	distribution := NewCapitalDistribution()
	for symbol, units := range a.Stocks {
		distribution.SetStock(symbol, units*window.MeanPrice(symbol))
	}
	return &distribution
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
	portfolio := NewAbstractPortfolio(0.0)
	for symbol, allocation := range a.stocks {
		portfolio.Stocks[symbol] = allocation * value / prices[symbol]
	}
	return portfolio
}

// Return an abstract portfolio representing a capital distribtuion at the given window's mean
// price.
func (a *CapitalDistribution) ToAbstractPortfolioOnMeanPrice(
	window Window,
	value float64,
) AbstractPortfolio {
	portfolio := NewAbstractPortfolio(0.0)
	for symbol, allocation := range a.stocks {
		portfolio.Stocks[symbol] = allocation * value / window.MeanPrice(symbol)
	}
	return portfolio
}

// Expected relative change from window open to close.
func (a *CapitalDistribution) WindowPerformanceAtMeanPrice(window Window) float64 {
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
