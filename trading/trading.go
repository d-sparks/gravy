package trading

import (
	"log"

	"github.com/Clever/go-utils/stringset"
	"github.com/d-sparks/gravy/db"
)

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

// Returns the capital distribution of this portfolio at the gven closing prices.
func (a *AbstractPortfolio) ToCapitalDistributionOnClose(prices map[string]db.Prices) *CapitalDistribution {
	cd := NewCapitalDistribution()
	for ticker, shares := range a.Stocks {
		cd.SetStock(ticker, shares*prices[ticker].Close)
	}
	cd.SetCash(a.CashUSD)
	return &cd
}

// CapitalDistribution portfolio. This is a portfolio normalized by value such that the sum of stocks and cash is 1.
type CapitalDistribution struct {
	stocks  map[string]float64
	cashUSD float64
	total   float64

	NonZeroStocks stringset.StringSet
}

func NewCapitalDistribution() CapitalDistribution {
	return CapitalDistribution{stocks: map[string]float64{}, total: 0.0, NonZeroStocks: stringset.New()}
}

func NewUniformCapitalDistribution(tickers stringset.StringSet) *CapitalDistribution {
	distribution := NewCapitalDistribution()
	for ticker, _ := range tickers {
		distribution.SetStock(ticker, 1.0)
	}
	distribution.SetCash(1.0)
	return &distribution
}

// Ensures that the abstract portfolio is normalized to be a distribution.
func (a *CapitalDistribution) ensureDistribution() {
	if a.total == 0.0 {
		log.Fatalf("Tried to create distribution in CapitalDistribution before SetStock.")
	} else if a.total == 1.0 {
		return
	}
	for ticker, value := range a.stocks {
		a.stocks[ticker] = value / a.total
	}
	a.cashUSD /= a.total
	a.total = 1.0
}

// Sets a value for stock and updates the total.
func (a *CapitalDistribution) SetStock(ticker string, value float64) {
	a.total -= a.stocks[ticker]
	a.stocks[ticker] = value
	a.total += value

	if value == 0.0 {
		a.NonZeroStocks.Remove(ticker)
	} else {
		a.NonZeroStocks.Add(ticker)
	}
}

// Sets cash holding for a capital distribution.
func (a *CapitalDistribution) SetCash(value float64) {
	a.total -= a.cashUSD
	a.cashUSD = value
	a.total += value
}

// Gets stock (after ensuring the portfolio is a distribution).
func (a *CapitalDistribution) GetStock(ticker string) float64 {
	a.ensureDistribution()
	return a.stocks[ticker]
}

// Gets cash value (after ensuring the distribution is
func (a *CapitalDistribution) GetCashUSD() float64 {
	a.ensureDistribution()
	return a.cashUSD
}

// Converts to the equivalent abstract portfolio at given prices.
func (a *CapitalDistribution) ToAbstractPortfolioOnClose(prices map[string]db.Prices, value float64) *AbstractPortfolio {
	a.ensureDistribution()
	portfolio := NewAbstractPortfolio(0.0)
	for ticker, allocation := range a.stocks {
		portfolio.Stocks[ticker] = allocation * value / prices[ticker].Close
	}
	return &portfolio
}
