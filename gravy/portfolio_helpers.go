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
