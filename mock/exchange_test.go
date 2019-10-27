package mock_test

import (
	"testing"

	"github.com/d-sparks/gravy/mock"
	"github.com/d-sparks/gravy/trading"
)

func TestExchange(t *testing.T) {
	seed := 1000.00
	e := mock.NewExchange(seed)

	prices := trading.Window{
		Open: trading.Prices{
			"FOO": 9,
			"BAR": 96,
		},
		Close: trading.Prices{
			"FOO": 10,
			"BAR": 100,
		},
		High: trading.Prices{
			"FOO": 15,
			"BAR": 100,
		},
		Low: trading.Prices{
			"FOO": 9,
			"BAR": 94,
		},
	}
	e.SetPrices(prices)

	order := trading.Order{
		Type:   "stop",
		Ticker: "BAR",
		Volume: 5,
		Price:  90,
	}
	_ = e.SubmitOrder(order)

	t.Run("limit order is applied if it falls within cash on hand", func(t *testing.T) {
		order := trading.Order{
			Type:   "limit",
			Ticker: "FOO",
			Volume: 20,
			Price:  9,
		}
		_ = e.SubmitOrder(order)

		p := e.CurrentPortfolio()
		if exp, got := 800.0, p.CashUSD; exp != got {
			t.Errorf("expected cash available to be %v, got %v", exp, got)
		}

		if exp, got := 20, p.Stocks["FOO"]; exp != got {
			t.Errorf("expected to have %v shares, got %v", exp, got)
		}

	})

	t.Run("new stop order is if price drops below expected position", func(t *testing.T) {
		order := trading.Order{
			Type:   "stop",
			Ticker: "FOO",
			Volume: 15,
			Price:  8,
		}
		_ = e.SubmitOrder(order)

		prices := trading.Window{
			Open: trading.Prices{
				"FOO": 9,
				"BAR": 96,
			},
			Close: trading.Prices{
				"FOO": 5,
				"BAR": 100,
			},
			High: trading.Prices{
				"FOO": 15,
				"BAR": 100,
			},
			Low: trading.Prices{
				"FOO": 9,
				"BAR": 94,
			},
		}
		e.SetPrices(prices)

		p := e.CurrentPortfolio()
		if exp, got := 875.0, p.CashUSD; exp != got {
			t.Errorf("expected cash available to be %v, got %v", exp, got)
		}

		if exp, got := 5, p.Stocks["FOO"]; exp != got {
			t.Errorf("expected to have %v shares, got %v", exp, got)
		}

	})

}
