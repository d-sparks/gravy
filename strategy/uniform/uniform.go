package uniform

import (
	"fmt"
	"time"

	"github.com/d-sparks/gravy/db"
	"github.com/d-sparks/gravy/db/dailyprices"
	"github.com/d-sparks/gravy/signal"
	"github.com/d-sparks/gravy/strategy"
	"github.com/d-sparks/gravy/trading"
)

const Name = "uniform"

// Uniform strategy. Invest equally in all securities.
type Uniform struct {
	debug map[string]string
}

func New() *Uniform {
	return &Uniform{}
}

func (b *Uniform) Name() string {
	return Name
}

// Return a uniform capital distribution every day based on the previous close price.
func (b *Uniform) Run(
	date time.Time,
	stores map[string]db.Store,
	signals map[string]signal.Signal,
) (*strategy.StrategyOutput, error) {
	b.debug = map[string]string{}

	// Get yesterday's prices.
	dailyPricesData, err := stores[dailyprices.Name].Get(date)
	if err != nil {
		return nil, fmt.Errorf("Error getting daily prices in strategy `%s`: `%s`", Name, err.Error())
	}

	// Create and return a uniform capital distribution.
	strategyOutput := strategy.StrategyOutput{
		CapitalDistribution: trading.NewUniformCapitalDistribution(dailyPricesData.Tickers),
	}

	//fmt.Println("=============================")
	//for ticker, _ := range strategyOutput.CapitalDistribution.NonZeroStocks {
	//	fmt.Printf("%s:%f,", ticker, strategyOutput.CapitalDistribution.GetStock(ticker))
	//}

	return &strategyOutput, nil
}

// No data to output
func (b *Uniform) Headers() []string {
	return []string{}
}

// Return last run's debug.
func (b *Uniform) Debug() map[string]string {
	return b.debug
}
