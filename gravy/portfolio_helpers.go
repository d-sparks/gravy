package gravy

import (
	"math"

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

// InvestApproximatelyUniformly attempts to invest approximately uniformly.
func InvestApproximatelyUniformlyInTargets(
	algorithmID *supervisor_pb.AlgorithmId,
	portfolio *supervisor_pb.Portfolio,
	prices *dailyprices_pb.DailyData,
	targets map[string]struct{},
) (orders []*supervisor_pb.Order) {
	if len(targets) == 0 {
		return []*supervisor_pb.Order{}
	}

	totalLimitOfOrders := 0.0
	target := 0.99 * PortfolioValue(portfolio, prices) / float64(len(targets))
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
		if _, ok := targets[ticker]; !ok {
			continue
		}
		volume := math.Floor(target / stockPrices.GetClose())
		if volume == 0.0 {
			continue
		}
		limit := 1.01 * stockPrices.GetClose()
		orders = append(orders, &supervisor_pb.Order{
			AlgorithmId: algorithmID, Ticker: ticker, Volume: volume, Limit: limit,
		})
		targetInvestments[ticker] = volume * limit
		totalLimitOfOrders += volume * limit
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
		limit := 1.01 * nextPrice
		orders = append(orders, &supervisor_pb.Order{
			AlgorithmId: algorithmID, Ticker: nextTicker, Volume: 1.0, Limit: limit,
		})
		targetInvestments[nextTicker] += limit
		totalLimitOfOrders += 1.0 * limit
	}

	return
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
