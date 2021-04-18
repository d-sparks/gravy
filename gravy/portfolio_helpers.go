package gravy

import (
	"math"
	"sort"

	dailyprices_pb "github.com/d-sparks/gravy/data/dailyprices/proto"
	supervisor_pb "github.com/d-sparks/gravy/supervisor/proto"
)

// PortfolioValue returns the value at last closing prices.
func PortfolioValue(portfolio *supervisor_pb.Portfolio, prices *dailyprices_pb.DailyData) float64 {
	value := portfolio.GetUsd()
	for ticker, volume := range portfolio.GetStocks() {
		value += volume * prices.GetPrices()[ticker].GetClose()
	}
	return value
}

// SignificantAllocations returns a co-indexed list of stocks and the weight of that stock in the portfolio, sorted by
// weight descending. Only returns the top 20 allocations amongst those with weight at least 1%.
func SignificantAllocations(
	portfolio *supervisor_pb.Portfolio,
	prices *dailyprices_pb.DailyData,
	portfolioValue float64,
) ([]string, []float64) {
	onePercentStocks := []struct {
		ticker string
		weight float64
	}{}

	// Find all heavy stocks.
	for ticker, volume := range portfolio.GetStocks() {
		weight := volume * prices.GetPrices()[ticker].GetClose() / portfolioValue
		if weight < 0.01 {
			continue
		}
		onePercentStocks = append(onePercentStocks, struct {
			ticker string
			weight float64
		}{ticker: ticker, weight: weight})
	}

	// Sort by weight descending.
	sort.Slice(onePercentStocks, func(i int, j int) bool {
		return onePercentStocks[i].weight > onePercentStocks[j].weight
	})

	// Build output.
	tickers := []string{}
	weights := []float64{}
	for i := 0; i < len(onePercentStocks) && i < 20; i++ {
		tickers = append(tickers, onePercentStocks[i].ticker)
		weights = append(weights, onePercentStocks[i].weight)
	}

	return tickers, weights
}

// TargetUniformInvestment returns the target value per stock to achieve uniform investment.
func TargetUniformInvestment(portfolio *supervisor_pb.Portfolio, prices *dailyprices_pb.DailyData) float64 {
	return PortfolioValue(portfolio, prices) / float64(len(prices.GetPrices()))
}

// SellEverythingWithStop creates stop orders to sell the portfolio.
func SellEverythingWithStop(
	algorithmID *supervisor_pb.AlgorithmId,
	portfolio *supervisor_pb.Portfolio,
	stop func(ticker string) float64,
) (orders []*supervisor_pb.Order) {
	for ticker, volume := range portfolio.GetStocks() {
		if volume == 0.0 {
			continue
		}
		orders = append(orders, &supervisor_pb.Order{
			AlgorithmId: algorithmID,
			Ticker:      ticker,
			Volume:      -volume,
			Stop:        stop(ticker),
		})
	}
	return
}

// SellEverythingWithStopPercent produces Orders to sell the entire portfolio. Each order has a stop equal to
// stopPercent * prices.GetPrices()[ticker].GetClose().
func SellEverythingWithStopPercent(
	algorithmID *supervisor_pb.AlgorithmId,
	portfolio *supervisor_pb.Portfolio,
	prices *dailyprices_pb.DailyData,
	stopPercent float64,
) []*supervisor_pb.Order {
	return SellEverythingWithStop(
		algorithmID,
		portfolio,
		func(ticker string) float64 { return prices.GetPrices()[ticker].GetClose() * stopPercent },
	)
}

// SellEverythingMarketOrder produces market orders to sell the entire portfolio.
func SellEverythingMarketOrder(
	algorithmID *supervisor_pb.AlgorithmId,
	portfolio *supervisor_pb.Portfolio,
) []*supervisor_pb.Order {
	return SellEverythingWithStop(algorithmID, portfolio, func(ticker string) float64 { return 0.0 })

}

// InvestApproximatelyUniformlyInTargetsWithLimits attempts to invest approximately uniformly with given UP and DOWN
// buffers. It is expected that upLimit * downLimit <= 1.0, so that the inlined comment below will hold in general.
func InvestApproximatelyUniformlyInTargetsWithLimits(
	algorithmID *supervisor_pb.AlgorithmId,
	portfolio *supervisor_pb.Portfolio,
	prices *dailyprices_pb.DailyData,
	targets map[string]struct{},
	investmentLimit float64,
	upLimit float64,
	downLimit float64,
	ignoreExisting bool,
) (orders []*supervisor_pb.Order) {
	if len(targets) == 0 {
		return []*supervisor_pb.Order{}
	}
	if investmentLimit == 0.0 {
		investmentLimit = portfolio.GetUsd()
	}

	totalLimitOfOrders := 0.0
	target := downLimit * PortfolioValue(portfolio, prices) / float64(len(targets))
	targetInvestments := map[string]float64{}

	// Note: this will represent a total limit of
	//
	//   1.01 * \sum_stock floor(target / price) * price <=
	//   1.01 * \sum_stock target =
	//   1.01 * 0.99 * \sum_stocks portfolioValue / # stocks =
	//   1.01 * 0.99 * portfolioValue =
	//   0.9999 * portfolioValue
	//
	// Thus, the investment is safe.
	for ticker, stockPrices := range prices.GetPrices() {
		if _, ok := targets[ticker]; !ok || stockPrices.GetClose() < 1e-4 {
			continue
		}
		limit := upLimit * stockPrices.GetClose()
		existing := 0.0
		if !ignoreExisting {
			existing = portfolio.GetStocks()[ticker]
		}
		volume := math.Floor(math.Min(
			target/stockPrices.GetClose()-existing,     // Ideal target
			(investmentLimit-totalLimitOfOrders)/limit, // Imposed limit
		))
		if volume <= 0.0 {
			continue
		}
		orders = append(orders, &supervisor_pb.Order{
			AlgorithmId: algorithmID, Ticker: ticker, Volume: volume, Limit: limit,
		})
		targetInvestments[ticker] = volume * limit
		totalLimitOfOrders += volume * limit // Because of min above, this is <= investmentLimit
	}

	// Continue investing until we can't get closer to a uniform investment.
	for portfolio.GetUsd()-totalLimitOfOrders > 0 {
		var nextTicker string
		nextImprovement := 0.0
		nextPrice := 0.0
		for ticker, stockPrices := range prices.GetPrices() {
			if _, ok := targets[ticker]; !ok {
				continue
			}
			currentTarget := targetInvestments[ticker]
			closePrice := stockPrices.GetClose()
			if closePrice+totalLimitOfOrders > portfolio.GetUsd() {
				continue
			}
			hypotheticalDelta := math.Abs(closePrice + currentTarget - target)
			currentDelta := math.Abs(currentTarget - target)
			improvement := currentDelta - hypotheticalDelta
			if improvement > nextImprovement {
				nextImprovement = improvement
				nextTicker = ticker
				nextPrice = closePrice
			}
		}
		if nextImprovement == 0.0 {
			// No improvement can be made.
			break
		}
		// Place an order for a nextTicker
		limit := upLimit * nextPrice
		orders = append(orders, &supervisor_pb.Order{
			AlgorithmId: algorithmID, Ticker: nextTicker, Volume: 1.0, Limit: limit,
		})
		targetInvestments[nextTicker] += limit
		totalLimitOfOrders += 1.0 * limit
	}

	return
}

// InvestApproximatelyUniformly attempts to invest approximately uniformly.
func InvestApproximatelyUniformlyInTargets(
	algorithmID *supervisor_pb.AlgorithmId,
	portfolio *supervisor_pb.Portfolio,
	prices *dailyprices_pb.DailyData,
	targets map[string]struct{},
) (orders []*supervisor_pb.Order) {
	return InvestApproximatelyUniformlyInTargetsWithLimits(
		algorithmID,
		portfolio,
		prices,
		targets,
		/*investmentLimit=*/ 0.0,
		/*upLimit=*/ 1.01,
		/*downLimit=*/ 0.99,
		/*ignoreExisting=*/ false,
	)
}

// InvestApproximatelyUniformly attempts to invest approximately uniformly.
func InvestApproximatelyUniformly(
	algorithmID *supervisor_pb.AlgorithmId,
	portfolio *supervisor_pb.Portfolio,
	prices *dailyprices_pb.DailyData,
) (orders []*supervisor_pb.Order) {
	targets := map[string]struct{}{}
	for ticker := range prices.GetPrices() {
		targets[ticker] = struct{}{}
	}
	return InvestApproximatelyUniformlyInTargets(algorithmID, portfolio, prices, targets)
}

// InvestApproximatelyUniformlyWithLimits invests in all targets.
func InvestApproximatelyUniformlyWithLimits(
	algorithmID *supervisor_pb.AlgorithmId,
	portfolio *supervisor_pb.Portfolio,
	prices *dailyprices_pb.DailyData,
	investmentLimit float64,
	upLimit float64,
	downLimit float64,
	ignoreExisting bool,
) (orders []*supervisor_pb.Order) {
	targets := map[string]struct{}{}
	for ticker := range prices.GetPrices() {
		targets[ticker] = struct{}{}
	}
	return InvestApproximatelyUniformlyInTargetsWithLimits(
		algorithmID,
		portfolio,
		prices,
		targets,
		investmentLimit,
		upLimit,
		downLimit,
		ignoreExisting,
	)
}
