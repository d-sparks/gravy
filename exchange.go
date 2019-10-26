package gravy

import "github.com/d-sparks/gravy/trading"

// Exchange is the interface to an exchange, e.g. a mock exchange or an alpaca client.
type Exchange interface {
	CurrentPortfolio() trading.Portfolio
	SubmitOrder(order trading.Order) trading.OrderOutcome
}

// Exchange for simulation.
type MockExchange struct {
	portfolio trading.Portfolio
	prices    trading.Window
}

// New Mock exchange starting with a seed of USD.
func NewMockExchange(seed float64) *MockExchange {
	return &MockExchange{portfolio: trading.NewPortfolio(seed)}
}

// Sets prices for upcoming orders.
func (m *MockExchange) SetPrices(prices trading.Window) {
	m.prices = prices
}

// Returns current portfolio.
func (m *MockExchange) CurrentPortfolio() trading.Portfolio {
	return m.portfolio
}

// Simulates an order based on a trading window during which it was placed. Must call SetPrices
// first. Updates portfolio accordingly.
func (m *MockExchange) SubmitOrder(order trading.Order) trading.OrderOutcome {
	return trading.OrderOutcome{}
}
