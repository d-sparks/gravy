package algorithm

import (
	"fmt"
	"time"

	"github.com/Clever/go-utils/stringset"
	"github.com/d-sparks/gravy"
	"github.com/d-sparks/gravy/db"
	"github.com/d-sparks/gravy/signal"
	"github.com/d-sparks/gravy/signal/ipos"
	"github.com/d-sparks/gravy/signal/movingaverage"
	"github.com/d-sparks/gravy/signal/unlistings"
	"github.com/d-sparks/gravy/strategy"
	"github.com/d-sparks/gravy/strategy/buyandhold"
	"github.com/d-sparks/gravy/trading"
)

// A TradingAlgorithm orchestrates stores, signals, strategies, and conducts trades.
type TradingAlgorithm struct {
	// For executing strategies and making trades.
	stores     map[string]db.Store
	signals    map[string]signal.Signal
	strategies map[string]strategy.Strategy
	exchange   gravy.Exchange

	// For debug.
	signalOrder      []string
	strategyOrder    []string
	nonhiddenHeaders stringset.StringSet
	headers          []string
	algorithmHeaders []string
	debug            map[string]string
}

func NewTradingAlgorithm(stores map[string]db.Store, exchange gravy.Exchange) TradingAlgorithm {
	// Initialize all members.
	algorithm := TradingAlgorithm{
		stores:     stores,
		signals:    map[string]signal.Signal{},
		strategies: map[string]strategy.Strategy{},
		exchange:   exchange,

		nonhiddenHeaders: stringset.New(),
		signalOrder:      []string{},
		strategyOrder:    []string{},
		debug:            map[string]string{},
	}

	// Initialize signals.
	algorithm.AddSignal(movingaverage.Name(100), movingaverage.New(100))
	algorithm.AddSignal(ipos.Name, ipos.New())
	algorithm.AddSignal(unlistings.Name, unlistings.New())

	// Initialize strategies.
	algorithm.AddStrategy(buyandhold.Name(100), buyandhold.New(100))

	// Order of algorithm headers. Use internal name (don't include algHeaders).
	algorithm.algorithmHeaders = []string{"date"}

	// Whitelist any nonhidden headers. Should match actual csv header, so use algHeader,
	// stratHeader, and signalHeader.
	algorithm.nonhiddenHeaders.Add(algHeader("date"))

	return algorithm
}

// Convenience for adding signals, since we also track the order they were added (to keep CSV
// header columns in order when printing debug).
func (t *TradingAlgorithm) AddSignal(name string, signal signal.Signal) {
	t.signals[name] = signal
	t.signalOrder = append(t.signalOrder, name)
}

// Convenience for adding strategies, since we also track the order they were added (to keep CSV
// header columns in order when printing debug).
func (t *TradingAlgorithm) AddStrategy(name string, strategy strategy.Strategy) {
	t.strategies[name] = strategy
	t.strategyOrder = append(t.strategyOrder, name)
}

// Calculates data, signals, strategies, and executes trades.
func (t *TradingAlgorithm) Trade(date time.Time) {
	// Clear debug output.
	t.debug = map[string]string{"date": date.Format("2006-01-02")}

	// Get current portfolio.
	// portfolio := t.exchange.CurrentPortfolio()

	// Get outputs of individual strategies.
	strategyOutputs := map[string]strategy.StrategyOutput{}
	for name, strategy := range t.strategies {
		strategyOutputs[name] = strategy.Run(date, t.stores, t.signals)
	}

	// Calculate orders
	orders := []trading.Order{}
	orderOutcomes := make([]trading.OrderOutcome, len(orders))
	for i, order := range orders {
		orderOutcomes[i] = t.exchange.SubmitOrder(order)
	}
}

// Format helpers for debug headers.
func signalHeader(signal, header string) string {
	return fmt.Sprintf("signalstrat.%s.%s", signal, header)
}
func stratHeader(strat, header string) string { return fmt.Sprintf("strat.%s.%s", strat, header) }
func algHeader(header string) string          { return fmt.Sprintf("alg.%s", header) }

// Combines all signal, strategy, and algorithm headers. On first call, actually computes header
// order.
func (t *TradingAlgorithm) Headers() []string {
	if len(t.headers) == 0 {
		t.headers = []string{}
		for _, signal := range t.signalOrder {
			for _, header := range t.signals[signal].Headers() {
				t.headers = append(t.headers, signalHeader(signal, header))
			}
		}
		for _, strat := range t.strategyOrder {
			for _, header := range t.strategies[strat].Headers() {
				t.headers = append(t.headers, stratHeader(strat, header))
			}
		}
		for _, header := range t.algorithmHeaders {
			t.headers = append(t.headers, algHeader(header))
		}
	}
	return t.headers
}

// Keyed by signalHeaders, stratHeaders, and algHeaders.
func (t *TradingAlgorithm) Debug(hide bool) map[string]string {
	debug := map[string]string{}

	// Get signal level debug.
	for name, signal := range t.signals {
		for header, value := range signal.Debug() {
			signalHeaderStr := signalHeader(name, header)
			if !hide || t.nonhiddenHeaders.Contains(signalHeaderStr) {
				debug[signalHeaderStr] = value
			}
		}
	}

	// Get strategy level debug.
	for name, strategy := range t.strategies {
		for header, value := range strategy.Debug() {
			stratHeaderStr := stratHeader(name, header)
			if !hide || t.nonhiddenHeaders.Contains(stratHeaderStr) {
				debug[stratHeaderStr] = value
			}
		}
	}

	// Get algorithm level debug.
	for header, value := range t.debug {
		algHeaderStr := algHeader(header)
		if !hide || t.nonhiddenHeaders.Contains(algHeaderStr) {
			debug[algHeaderStr] = value
		}
	}

	return debug
}
