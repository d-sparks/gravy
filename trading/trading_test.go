package trading

import (
	"testing"

	"github.com/Clever/go-utils/stringset"
	"github.com/stretchr/testify/assert"
)

var window Window

func init() {
	window = Window{
		Open:    Prices{"MSFT": 1.5, "GOOGL": 2.5},
		Close:   Prices{"MSFT": 2.5, "GOOGL": 3.5},
		Low:     Prices{"MSFT": 1.0, "GOOGL": 2.0},
		High:    Prices{"MSFT": 3.0, "GOOGL": 4.0},
		Symbols: stringset.New("MSFT", "GOOGL"),
	}
}

func TestWindowMeanHighLowPrice(t *testing.T) {
	assert.Equal(t, 2.0, window.MeanHighLowPrice("MSFT"))
	assert.Equal(t, 3.0, window.MeanHighLowPrice("GOOGL"))
}

func TestPortfolioValue(t *testing.T) {
	portfolio := NewPortfolio(1.0)
	portfolio.Stocks["MSFT"] = 2
	portfolio.Stocks["GOOGL"] = 3
	assert.Equal(t, 19.0, portfolio.Value(window.High))
}

func TestPortfolioMeanHighLowValue(t *testing.T) {
	portfolio := NewPortfolio(1.0)
	portfolio.Stocks["MSFT"] = 2
	portfolio.Stocks["GOOGL"] = 3
	assert.Equal(t, 14.0, portfolio.MeanHighLowValue(window))
}

func TestAbstractPortfolioValue(t *testing.T) {
	portfolio := NewAbstractPortfolio(1.0)
	portfolio.Stocks["MSFT"] = 2.0
	portfolio.Stocks["GOOGL"] = 3.0
	assert.Equal(t, 19.0, portfolio.Value(window.High))
}

func TestAbstractPortfolioMeanHighLowValue(t *testing.T) {
	portfolio := NewAbstractPortfolio(1.0)
	portfolio.Stocks["MSFT"] = 2.0
	portfolio.Stocks["GOOGL"] = 3.0
	assert.Equal(t, 14.0, portfolio.MeanHighLowValue(window))
}

func TestNewBalancedCapitalDistribution(t *testing.T) {
	distribution := NewBalancedCapitalDistribution(window.High)
	assert.Equal(t, distribution.GetStock("MSFT"), 0.5)
	assert.Equal(t, distribution.GetStock("GOOGL"), 0.5)
}

func TestCapitalDistributionGetSetStock(t *testing.T) {
	distribution := NewCapitalDistribution()
	distribution.SetStock("gravy", 1.0)
	assert.Equal(t, 1.0, distribution.GetStock("gravy"))
}

func TestCapitalDistributionToAbstractPortfolioOnPrices(t *testing.T) {
	distribution := NewBalancedCapitalDistribution(window.High)
	portfolio := distribution.ToAbstractPortfolioOnPrices(window.High, 2.0)
	assert.Equal(t, 1.0, window.High["MSFT"]*portfolio.Stocks["MSFT"])
	assert.Equal(t, 1.0, window.High["GOOGL"]*portfolio.Stocks["GOOGL"])
}

func TestCapitalDistributionRelativeWindowPerformance(t *testing.T) {
	// 4 shares msft, 3 shares googl
	distribution := NewBalancedCapitalDistribution(window.Open)
	total := (0.5)*(2.5-1.5)/1.5 + (0.5)*(3.5-2.5)/2.5
	assert.Equal(t, total, distribution.RelativeWindowPerformance(window))
}
