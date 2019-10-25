package algorithm

import (
	"time"

	"github.com/Clever/go-utils/stringset"
	"github.com/d-sparks/gravy/db"
	"github.com/d-sparks/gravy/exchange"
	"github.com/d-sparks/gravy/signal"
	"github.com/d-sparks/gravy/strategy"
	"github.com/d-sparks/gravy/trading"
)

// A TradingAlgorithm orchestrates stores, signals, strategies, and conducts trades.
type TradingAlgorithm struct {
	// For executing strategies and making trades.
	stores     map[string]db.Store
	signals    map[string]signal.Signal
	strategies map[string]strategy.Strategy
	exchange   exchange.Exchange

	// For debug.
	signalOrder      []string
	strategyOrder    []string
	nonhiddenHeaders stringset.StringSet
	headers          []string
	algorithmHeaders []string
	debug            map[string]string
}

func NewTradingAlgorithm(stores map[string]db.Store, exchange exchange.Exchange) TradingAlgorithm {
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
	algorithm.signals[movingaverage.Name(100)] = NewMovingAverage(100)
	algorithm.signalOrder = append(algorithm.signalOrder, movingaverage.Name(100))

	// Initialize strategies.

	// Order of algorithm headers.
	algorithm.algorithmHeaders = []string{"date"}

	// Whitelist anynonhidden headers.
	nonhiddenHeaders.Insert("date")
}

func (t *TradingAlgorithm) Trade(date time.Time) {
	// Clear debug output.
	t.debug = map[string]string{}

	// Get current portfolio.
	portfolio := s.exchange.CurrentPortfolio()

	// Get outputs of individual strategies.
	strategyOutputs := map[string]strategy.StrategyOutput{}
	for name, strategy := range t.strategies {
		strategyOutputs[name] = strategy.Run(date, t.stores, t.signals)
	}

	// Calculate orders
	orders := []trading.Order{}
	orderOutcomes := make([]trading.OrderOutcome, len(orders))
	for i, order := range orders {
		orderOutcomes[i] = exchange.SubmitOrder(order)
	}
}

func signalHeader(signal, header string) { return fmt.Sprintf("signalstrat.%s.%s", signal, header) }

func stratHeader(strat, header string) { return fmt.Sprintf("strat.%s.%s", strat, header) }

func algHeader(header string) { return fmt.Spritnf("alg.%s", header) }

// Returns headers signal.${SIGNAL}.header, strategy.${STRATEGY}.header, and its own headers. Must
// be called after signals and strategies are populated.
func (t *TradingAlgorithm) Headers() []string {
	// If headers has not been called, compute the header vector.
	if len(t.headers) == 0 {
		t.headers = []string{}
		for _, signal := range t.signalOrder {
			for _, header := range t.signals[signal].DebugHeaders() {
				t.headers = append(t.headers, signalHeader(signal, header))
			}
		}
		for _, strat := range t.strategyOrder {
			for _, header := range t.signals[strat].DebugHeaders() {
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
