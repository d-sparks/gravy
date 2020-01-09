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

const Name = "buyandhold"

// BuyAndHold strategy. Invest equally in all securities. When there is an IPO or unlisting.
type BuyAndHold struct {
	abstractPortfolio *trading.AbstractPortfolio
	debug             map[string]string
}

func New() *BuyAndHold {
	return &BuyAndHold{}
}

func (b *BuyAndHold) Name() string {
	return Name
}

// Whenever there is an ipo or unlisting, invest equally in all securities. This is achieved by taking a uniform
// capital distribution and turning it into an abstract portfolio worth $1.0 on the most recent closing prices. This
// portfolio is memorized until the next ipo or unlisting. On any other day, we imagine having held that original
// abstract portfolio and recommend a capital distribution that replicates that portfolio at the most recent prices.
func (b *BuyAndHold) Run(
	date time.Time,
	stores map[string]db.Store,
	signals map[string]signal.Signal,
) (*strategy.StrategyOutput, error) {
	b.debug = map[string]string{}

	iposData, err := signals[ipos.Name].Compute(date, stores)
	if err != nil {
		return nil, fmt.Errorf("Error computing ipos in strategy `%s`: `%s`", Name, err.Error())
	}
	unlistingsData, err := signals[unlistings.Name].Compute(date, stores)
	if err != nil {
		return nil, fmt.Errorf("Error computing unlistings in strategy `%s`: `%s`", Name, err.Error())
	}
	dailyPricesData, err := stores[dailyprices.Name].Get(date)
	if err != nil {
		return nil, fmt.Errorf("Error getting daily prices in strategy `%s`: `%s`", Name, err.Error())
	}

	if b.abstractPortfolio == nil || len(iposData.StringSet) > 0 || len(unlistingsData.StringSet) > 0 {
		uniform := trading.NewUniformCapitalDistribution(dailyPricesData.Tickers)
		b.abstractPortfolio = uniform.ToAbstractPortfolioOnClose(dailyPricesData.TickersToPrices, 1.0)
	}

	strategyOutput := strategy.StrategyOutput{
		CapitalDistribution: b.abstractPortfolio.ToCapitalDistributionOnClose(dailyPricesData.TickersToPrices),
	}
	return &strategyOutput, nil
}

// No data to output
func (b *BuyAndHold) Headers() []string {
	return []string{}
}

// Return last run's debug.
func (b *BuyAndHold) Debug() map[string]string {
	return b.debug
}
