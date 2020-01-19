package buyandhold

import (
	"fmt"
	"time"

	"github.com/d-sparks/gravy/db"
	"github.com/d-sparks/gravy/db/dailyprices"
	"github.com/d-sparks/gravy/signal"
	"github.com/d-sparks/gravy/signal/ipos"
	"github.com/d-sparks/gravy/signal/unlistings"
	"github.com/d-sparks/gravy/strategy"
	"github.com/d-sparks/gravy/trading"
)

// Name of the buyandhold strategy.
const Name = "buyandhold"

// BuyAndHold strategy. Invest equally in all securities. When there is an IPO or unlisting.
type BuyAndHold struct {
	abstractPortfolio *trading.AbstractPortfolio
	previousIPOs      bool
	debug             map[string]string
}

// New BuyAndHold strategy.
func New() *BuyAndHold {
	return &BuyAndHold{abstractPortfolio: nil, previousIPOs: false, debug: map[string]string{}}
}

// Name returns the name as per the strategy interface.
func (b *BuyAndHold) Name() string {
	return Name
}

// Run the strategy. In the first run, when there is an unlisting, or the session after an IPO signal, create an
// abstract portfolio invested uniformly in (tomorrow's tickers) \bigcap (today's tickers) based on today's prices.
// "Hold" this abstract portfolio by returning a capital distribution that is the distribution of the mature value of
// the portfolio.
func (b *BuyAndHold) Run(
	date time.Time,
	stores map[string]db.Store,
	signals map[string]signal.Signal,
) (*strategy.StrategyOutput, error) {
	b.debug = map[string]string{}

	// Get unlistings and prices.
	unlistingsData, err := signals[unlistings.Name].Compute(date, stores)
	if err != nil {
		return nil, fmt.Errorf("Error computing unlistings in strategy `%s`: `%s`", Name, err.Error())
	}
	pricesData, err := stores[dailyprices.Name].Get(date)
	if err != nil {
		return nil, fmt.Errorf("Error getting daily prices in strategy `%s`: `%s`", Name, err.Error())
	}

	// If we need to rebalance our portfolio, do so.
	if b.abstractPortfolio == nil || b.previousIPOs || len(unlistingsData.StringSet) > 0 {
		investIn := pricesData.Tickers
		for ticker := range unlistingsData.StringSet {
			investIn.Remove(ticker)
		}
		uniform := trading.NewUniformCapitalDistribution(investIn)
		b.abstractPortfolio = uniform.ToAbstractPortfolioOnClose(pricesData.TickersToPrices, 1.0)
	}

	// If there were IPOs, remember this so we can rebalance in the next session.
	iposData, err := signals[ipos.Name].Compute(date, stores)
	if err != nil {
		return nil, fmt.Errorf("Error computing ipos in strategy `%s`: `%s`", Name, err.Error())
	}
	b.previousIPOs = len(iposData.StringSet) > 0

	// Return a capital distribution corresponding to the mature value of the abstract portfolio.
	strategyOutput := strategy.StrategyOutput{
		CapitalDistribution: b.abstractPortfolio.ToCapitalDistributionOnClose(pricesData.TickersToPrices),
	}
	return &strategyOutput, nil
}

// Headers as per the strategy interface.
func (b *BuyAndHold) Headers() []string {
	return []string{}
}

// Debug as per the strategy interface.
func (b *BuyAndHold) Debug() map[string]string {
	return b.debug
}
