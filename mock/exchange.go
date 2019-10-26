package mock

import "github.com/d-sparks/gravy/trading"

// Exchange for simulation.
type Exchange struct {
	portfolio trading.Portfolio
	prices    trading.Window
}

// New exchange starting with a seed of USD.
func NewExchange(seed float64) *Exchange {
	return &Exchange{portfolio: trading.NewPortfolio(seed)}
}

// Sets prices for upcoming orders.
func (m *Exchange) SetPrices(prices trading.Window) {
	m.prices = prices
}

// Returns current portfolio.
func (m *Exchange) CurrentPortfolio() trading.Portfolio {
	return m.portfolio
}

// Simulates an order based on a trading window during which it was placed. Must call SetPrices
// first. Updates portfolio accordingly.
func (m *Exchange) SubmitOrder(order trading.Order) trading.OrderOutcome {
	return trading.OrderOutcome{}
}
