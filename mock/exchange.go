package mock

import (
	"fmt"
	"sync"

	"github.com/d-sparks/gravy/trading"
)

// Exchange for simulation.
type Exchange struct {
	mu        sync.Mutex
	portfolio trading.Portfolio
	prices    trading.Window

	orders []trading.Order
}

// New exchange starting with a seed of USD.
func NewExchange(seed float64) *Exchange {
	return &Exchange{portfolio: trading.NewPortfolio(seed)}
}

// Sets prices for upcoming orders.
func (m *Exchange) SetPrices(prices trading.Window) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.prices = prices

	m.applyOrders()
}

// Returns current portfolio.
func (m *Exchange) CurrentPortfolio() trading.Portfolio {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.portfolio
}

// Simulates an order based on a trading window during which it was placed. Must call SetPrices
// first. Updates portfolio accordingly.
func (m *Exchange) SubmitOrder(order trading.Order) trading.OrderOutcome {
	m.mu.Lock()
	defer m.mu.Unlock()

	switch order.Type {
	case "limit", "stop":
		m.orders = append(m.orders, order)
	default:
		panic(fmt.Errorf("uknown order type %s", order.Type))
	}

	m.applyOrders()

	return trading.OrderOutcome{}
}

func (m *Exchange) applyOrders() {
	unfilledOrders := []trading.Order{}
	for _, order := range m.orders {
		ticker := order.Ticker
		price := m.prices.Close[ticker]
		orderCost := price * float64(order.Volume)

		switch order.Type {
		case "limit":
			if price >= order.Price && orderCost <= m.portfolio.CashUSD {
				m.portfolio.Stocks[ticker] += order.Volume
				m.portfolio.CashUSD -= orderCost
			} else {
				unfilledOrders = append(unfilledOrders, order)
			}
		case "stop":
			if price <= order.Price {
				// This potentially allows users to end up in a short position
				m.portfolio.Stocks[ticker] -= order.Volume
				m.portfolio.CashUSD += orderCost
			} else {
				unfilledOrders = append(unfilledOrders, order)
			}
		}
	}

	m.orders = unfilledOrders
}
