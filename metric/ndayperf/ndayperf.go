package ndayperf

import (
	"fmt"
	"time"

	"github.com/d-sparks/gravy/db"
	"github.com/d-sparks/gravy/db/dailyprices"
	"github.com/d-sparks/gravy/signal"
	"github.com/d-sparks/gravy/strategy"
)

func Name(days int) string { return fmt.Sprintf("%ddayperf", days) }

// Performance over N days.
type NDayPerf struct {
	days int
}

func New(days int) *NDayPerf {
	return &NDayPerf{days}
}

// Return the name for this metric.
func (n *NDayPerf) Name() string {
	return Name(n.days)
}

func (n *NDayPerf) Value(
	date time.Time,
	stores map[string]db.Store,
	signals map[string]signal.Signal,
	strategyOutput *strategy.StrategyOutput,
) (float64, error) {
	// Get today's prices.
	todayPricesData, err := stores[dailyprices.Name].Get(date)
	if err != nil {
		fmt.Errorf("Error getting todays prices in  `%s`: `%s`", n.Name(), err.Error())
	}

	// Get prices after n days.
	endDate := date.AddDate(0, 0, n.days)
	endPricesData, err := stores[dailyprices.Name].Get(endDate)
	if err != nil {
		fmt.Errorf("Error getting end prices in `%s`: `%s`", n.Name(), err.Error())
	}

	// Calculate rate of return of each stock in the portfolio.
	cd := strategyOutput.CapitalDistribution
	overallReturn := 0.0
	for ticker, _ := range cd.NonZeroStocks {
		todayPrices, ok := todayPricesData.TickersToPrices[ticker]
		if !ok {
			return 0.0, fmt.Errorf("No opening price for `%s` in metric `%s`: `%s`", ticker, n.Name())
		}

		endPrices, ok := endPricesData.TickersToPrices[ticker]
		if ok {
			overallReturn += cd.GetStock(ticker) * (endPrices.Close / todayPrices.Close)
		} else {
			// If there is no closing price assume a neutral return. (This might be a bad assumption.)
			overallReturn += cd.GetStock(ticker) // * 1.0
		}
	}

	return overallReturn, nil
}
