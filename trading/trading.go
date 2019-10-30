package trading

import (
	"log"
	"time"

	"github.com/Clever/go-utils/stringset"
)

// Prices is a mapping from asset identifier to a price.
type Prices map[string]float64

// Window represents a unit of time, with various prices for the assets.
type Window struct {
	Begin time.Time
	End   time.Time

	Close  Prices
	High   Prices
	Low    Prices
	Open   Prices
	Volume Prices

	Symbols stringset.StringSet
}

func (w Window) MeanHighLowPrice(symbol string) float64 {
	return (w.High[symbol] + w.Low[symbol]) / 2.0
}

// Portfolio.
type Portfolio struct {
	Stocks  map[string]int
	CashUSD float64
}

func NewPortfolio(seed float64) Portfolio {
	return Portfolio{
		Stocks:  map[string]int{},
		CashUSD: seed,
	}

}

// Returns mature value of a position given a Tick.
func (p *Portfolio) Value(prices Prices) float64 {
	value := p.CashUSD
	for symbol, units := range p.Stocks {
		value += float64(units) * prices[symbol]
	}
	return value
}

// Mean value during a window.
func (p *Portfolio) MeanHighLowValue(window Window) float64 {
	value := p.CashUSD
	c := 0.0
	for symbol, units := range p.Stocks {
		// value += window.Close[symbol] * units //, with Kahan summation
		v := window.MeanHighLowPrice(symbol)*float64(units) - c
		t := value + v
		c = (t - value) - v
		value = t
	}
	return value
}

// AbstractPortfolio. Can have fractional/negative shares.
type AbstractPortfolio struct {
	Stocks  map[string]float64
	CashUSD float64
}

func NewAbstractPortfolio(seed float64) AbstractPortfolio {
	return AbstractPortfolio{
		Stocks:  map[string]float64{},
		CashUSD: seed,
	}
}

// Returns mature value of a position given a Tick.
func (a *AbstractPortfolio) Value(prices Prices) float64 {
	value := a.CashUSD
	for symbol, units := range a.Stocks {
		value += float64(units) * prices[symbol]
	}
	return value
}

// Mean value during a window.
func (a *AbstractPortfolio) MeanHighLowValue(window Window) float64 {
	value := a.CashUSD
	c := 0.0
	for symbol, units := range a.Stocks {
		// value += window.Close[symbol] * units //, with Kahan summation
		v := window.MeanHighLowPrice(symbol)*units - c
		t := value + v
		c = (t - value) - v
		value = t
	}
	return value
}

// CapitalDistribution portfolio. A probability distribution over potential assets.
type CapitalDistribution struct {
	stocks map[string]float64
	total  float64
}

func NewCapitalDistribution() CapitalDistribution {
	return CapitalDistribution{stocks: map[string]float64{}, total: 0.0}
}

func NewBalancedCapitalDistribution(prices Prices) CapitalDistribution {
	distribution := NewCapitalDistribution()
	for symbol, _ := range prices {
		distribution.SetStock(symbol, 1.0)
	}
	return distribution
}

// Ensures that the abstract portfolio is normalized to be a distribution.
func (a *CapitalDistribution) ensureDistribution() {
	if a.total == 0.0 {
		log.Fatalf("Tried to create distribution in CapitalDistribution before SetStock.")
	} else if a.total == 1.0 {
		return
	}
	for symbol, value := range a.stocks {
		a.stocks[symbol] = value / a.total
	}
	a.total = 1.0
}

// Sets a value for stock and updates the total.
func (a *CapitalDistribution) SetStock(symbol string, value float64) {
	a.total -= a.stocks[symbol]
	a.stocks[symbol] = value
	a.total += value
}

// Gets stock (after ensuring the portfolio is a distribution).
func (a *CapitalDistribution) GetStock(symbol string) float64 {
	a.ensureDistribution()
	return a.stocks[symbol]
}

// Converts to the equivalent abstract portfolio at given prices.
func (a *CapitalDistribution) ToAbstractPortfolioOnPrices(
	prices Prices,
	value float64,
) AbstractPortfolio {
	a.ensureDistribution()
	portfolio := NewAbstractPortfolio(0.0)
	for symbol, allocation := range a.stocks {
		portfolio.Stocks[symbol] = allocation * value / prices[symbol]
	}
	return portfolio
}

// Expected relative change from window open to close.
func (a *CapitalDistribution) RelativeWindowPerformance(window Window) float64 {
	a.ensureDistribution()
	perf := 0.0
	for symbol, exposure := range a.stocks {
		if exposure == 0.0 {
			continue
		} else if !window.Symbols.Contains(symbol) {
			log.Fatalf("Capital in unlisted symbol")
		}
		open := window.Open[symbol]
		perf += exposure * (window.Close[symbol] - open) / open
	}
	return perf
}

// OrderType is the type of order that will be submitted. The definitions
// for each type can be found https://www.thebalance.com/making-sense-of-day-trading-order-types-1031387
type OrderType int8

const (
	// This ensures that the zero value is not a valid OrderType
	UnknownOrderType OrderType = iota

	// BuyMarketOrderType buys you whatever price is available in the marketplace. The problem
	// with market orders is that you don't know the exact price you will end up buying at
	BuyMarketOrderType = iota + 1

	// SellMarketOrderType sells you whatever price is available in the marketplace. The problem
	// with market orders is that you don't know the exact price you will end up selling at
	SellMarketOrderType

	// BuyLimitOrderType is an order that is placed below the current price. The order
	// is only filled at or below the limit price. Buy limits are also used as targets,
	// to get you out of a profitable short trade.
	BuyLimitOrderType

	// SellLimitOrderType is an order to sell that is placed above the current price. The order
	// is only filled at or above the limit price. Sell limits are also often used as targets,
	// to get you out of a profitable long trade.
	SellLimitOrderType

	// BuyStopOrderType is an order to buy that is place above the current price. The order is only
	// is only filled at or above the stop price. Buy stops act like market orders once the buy stop
	// price is reached. They are also useful for buying breakouts above resistance, but you can't be
	// sure of the exact price you will end up buying at. Therefore, they are useful for using as a stop
	// loss on short positions, when you must get out because the price is moving against you.
	BuyStopOrderType

	// SellStopOrderType is an order to sell that is placed below the current price. The order will only
	// be filled at or below the stop price. This order can be used to get out of a long trade. Sell stops
	// act like market orders once the sell stop price is reached. Therefore, they are useful for using as a
	// stop loss on long positions, when you must get out because the price is moving against you.
	SellStopOrderType

	// BuyStopLimitOrder is very similar to a Buy Stop order, except that it doesn't act like a market order. The
	// buy stop limit will only fill at the buy stop limit price or lower. A buy stop limit order is useful for buying
	// when the price breaks above a particular level (such a resistance) but you only want to buy at a specific price
	// or lower when that event occurs.
	BuyStopLimitOrderType

	// SellStopLimitOrder is very similar to a Sell Stop order, except that it doesn't act like a market order. The sell
	// stop limit will only fill at the price equivalent to the limit price attached to the order, or higher. A sell stop
	// limit order is useful for selling when the price breaks below a particular level (such a support), but you only want
	// to sell at a specific price or higher when that event occurs.
	SellStopLimitOrderType
)

// String implements string.Stringer for OrderType.
func (t OrderType) String() string {
	switch t {
	case BuyMarketOrderType:
		return "buy_market"
	case SellMarketOrderType:
		return "sell_market"
	case BuyLimitOrderType:
		return "buy_limit"
	case SellLimitOrderType:
		return "sell_limit"
	case BuyStopOrderType:
		return "buy_stop"
	case SellStopOrderType:
		return "sell_stop"
	case BuyStopLimitOrderType:
		return "buy_stop_limit"
	case SellStopLimitOrderType:
		return "sell_stop_limit"
	default:
		return "unknown"
	}
}

// Order is a trading order. It specifies the volume, price, ticker and type.
type Order struct {
	Type   OrderType
	Ticker string
	Price  float64
	Volume int
}

// Not sure what we're going to want here.
type OrderOutcome struct {
}
