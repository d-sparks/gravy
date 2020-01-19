package trading

import (
	"log"

	"github.com/Clever/go-utils/stringset"
	"github.com/d-sparks/gravy/db"
)

// AbstractPortfolio is a portfolio that can have fractional/negative shares.
type AbstractPortfolio struct {
	Stocks  map[string]float64
	CashUSD float64
}

// NewAbstractPortfolio initializes an abstract portfolio with a given amount of seed capital in USD.
func NewAbstractPortfolio(seed float64) AbstractPortfolio {
	return AbstractPortfolio{
		Stocks:  map[string]float64{},
		CashUSD: seed,
	}
}

// ToCapitalDistributionOnClose returns the distribution of capital.
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

// NewCapitalDistribution returns an empty (invalid) capital distribution.
func NewCapitalDistribution() CapitalDistribution {
	return CapitalDistribution{stocks: map[string]float64{}, total: 0.0, NonZeroStocks: stringset.New()}
}

// NewUniformCapitalDistribution returns a capital distribution with no USD but equally invested in every equity.
func NewUniformCapitalDistribution(tickers stringset.StringSet) *CapitalDistribution {
	distribution := NewCapitalDistribution()
	for ticker := range tickers {
		distribution.SetStock(ticker, 1.0)
	}
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

// SetStock sets the value of a stock. All calls to SetStock should happen before a call to GetStock and keep in mind
// that calling SetStock after GetStock will mean the capital distribution was already normalized.
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

// SetCash sets the cash holdings and has the same caveat as SetStock.
func (a *CapitalDistribution) SetCash(value float64) {
	a.total -= a.cashUSD
	a.cashUSD = value
	a.total += value
}

// GetStock gets the stock value. This will cause the distribution to normalize.
func (a *CapitalDistribution) GetStock(ticker string) float64 {
	a.ensureDistribution()
	return a.stocks[ticker]
}

// GetCashUSD returns the cash USD probability. This will cause the distribution to normalize.
func (a *CapitalDistribution) GetCashUSD() float64 {
	a.ensureDistribution()
	return a.cashUSD
}

// ToAbstractPortfolioOnClose returns an abstract portfolio that is distributed as the capital distribution (at the
// given closing prices).
func (a *CapitalDistribution) ToAbstractPortfolioOnClose(prices map[string]db.Prices, value float64) *AbstractPortfolio {
	a.ensureDistribution()
	portfolio := NewAbstractPortfolio(0.0)
	for ticker, allocation := range a.stocks {
		portfolio.Stocks[ticker] = allocation * value / prices[ticker].Close
	}
	return &portfolio
}
