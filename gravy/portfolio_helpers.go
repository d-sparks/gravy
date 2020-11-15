package gravy

import (
	dailyprices_pb "github.com/d-sparks/gravy/data/dailyprices/proto"
	supervisor_pb "github.com/d-sparks/gravy/supervisor/proto"
)

// PortfolioValue returns the value at last closing prices.
func PortfolioValue(portfolio *supervisor_pb.Portfolio, prices *dailyprices_pb.DailyPrices) float64 {
	value := portfolio.GetUsd()
	for ticker, volume := range portfolio.GetStocks() {
		value += volume * prices.GetStockPrices()[ticker].GetClose()
	}
	return value
}

// TargetUniformInvestment returns the target value per stock to achieve uniform investment.
func TargetUniformInvestment(portfolio *supervisor_pb.Portfolio, prices *dailyprices_pb.DailyPrices) float64 {
	return PortfolioValue(portfolio, prices) / float64(len(prices.GetStockPrices()))
}

// SellEverythingWithStop creates stop orders to sell the portfolio.
func SellEverythingWithStop(
	algorithmID *supervisor_pb.AlgorithmId,
	portfolio *supervisor_pb.Portfolio,
	stop func(ticker string)float64,
) (orders []*supervisor_pb.Order) {
	for ticker, volume := range portfolio.GetStocks() {
		orders = append(orders, &supervisor_pb.Order{
			AlgorithmId: algorithmID,
			Ticker: ticker,
			Volume: -volume,
			Stop: stop(ticker),
		})
	}
	return
}

// SellEverythingWithStopPercent produces Orders to sell the entire portfolio. Each order has a stop equal to
// stopPercent * prices.GetStockPrices()[ticker].GetClose().
func SellEverythingWithStopPercent(
	algorithmID *supervisor_pb.AlgorithmId,
	portfolio *supervisor_pb.Portfolio,
	prices *dailyprices_pb.DailyPrices,
	stopPercent float64,
) []*supervisor_pb.Order {
	return SellEverythingWithStop(
		algorithmID,
		portfolio,
		func(ticker string) float64 { return prices.GetStockPrices()[ticker].GetClose() * stopPercent },
	)
}

// SellEverythingMarketOrder produces market orders to sell the entire portfolio.
func SellEverythingMarketOrder(
	algorithmID *supervisor_pb.AlgorithmId,
	portfolio *supervisor_pb.Portfolio,
) []*supervisor_pb.Order {
	return SellEverythingWithStop(algorithmID, portfolio, func(ticker string) float64 { return 0.0 })

}