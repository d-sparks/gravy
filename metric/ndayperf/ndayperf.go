package ndayperf

import (
	"fmt"
	"time"

	"github.com/d-sparks/gravy/db"
	"github.com/d-sparks/gravy/db/dailyprices"
	"github.com/d-sparks/gravy/signal"
	"github.com/d-sparks/gravy/strategy"
)

// Name of the ndayperf metric.
func Name(days int) string { return fmt.Sprintf("%ddayperf", days) }

// NDayPerf is the performance over N days.
type NDayPerf struct {
	days int
}

// New NDayPerf strategy.
func New(days int) *NDayPerf {
	return &NDayPerf{days}
}

// Name as per the metric interface.
func (n *NDayPerf) Name() string {
	return Name(n.days)
}

// Value of the metric is the rate of return of a capital distribution over N days. That is
//
//   \sum_ticker probability(ticker) * close(ticker, today + N) / open(ticker, today + 1)
//
// This is a percentage and can be positive or negative. The question of what to do when tickers from today + 1 are not
// still listed at today + N is left open. For now, for such stocks we assume a neutral rate of return.
func (n *NDayPerf) Value(
	date time.Time,
	stores map[string]db.Store,
	signals map[string]signal.Signal,
	strategyOutput *strategy.StrategyOutput,
) (float64, error) {
	// Get tomorrow's prices.
	nextTradingDay, err := stores[dailyprices.Name].NextDate(date)
	if err != nil {
		return 0.0, fmt.Errorf("Couldn't get nextTradingDay in `%s`: `%s`", n.Name(), err.Error())
	}
	tomorrowPricesData, err := stores[dailyprices.Name].Get(*nextTradingDay)
	if err != nil {
		return 0.0, fmt.Errorf("Error getting todays prices in  `%s`: `%s`", n.Name(), err.Error())
	}

	// Get prices after n days. The minus one is because NextDate will also increase by at least 1.
	endDate := date.AddDate(0, 0, n.days-1)
	endTradingDate, err := stores[dailyprices.Name].NextDate(endDate)
	if err != nil {
		return 0.0, fmt.Errorf("Couldn't get nDay trading date in `%s`: `%s`", n.Name(), err.Error())
	}
	endPricesData, err := stores[dailyprices.Name].Get(*endTradingDate)
	if err != nil {
		return 0.0, fmt.Errorf("Error getting end prices in `%s`: `%s`", n.Name(), err.Error())
	}

	// Calculate rate of return of each stock in the portfolio.
	cd := strategyOutput.CapitalDistribution
	overallReturn := 0.0
	for ticker := range cd.NonZeroStocks {
		tomorrowPrices, ok := tomorrowPricesData.TickersToPrices[ticker]
		if !ok {
			return 0.0, fmt.Errorf("No opening price for `%s` in metric `%s`", ticker, n.Name())
		}

		endPrices, ok := endPricesData.TickersToPrices[ticker]
		if ok {
			overallReturn += cd.GetStock(ticker) * (endPrices.Close / tomorrowPrices.Open)
		} else {
			// If there is no closing price assume a neutral return. (This might be a bad assumption.)
			overallReturn += cd.GetStock(ticker) // * 1.0
		}
	}

	return overallReturn, nil
}
