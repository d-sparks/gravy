package algorithm

import (
	"fmt"
	"strconv"
	"time"

	"github.com/Clever/go-utils/stringset"
	"github.com/d-sparks/gravy"
	"github.com/d-sparks/gravy/db"
	"github.com/d-sparks/gravy/metric"
	"github.com/d-sparks/gravy/metric/ndayperf"
	"github.com/d-sparks/gravy/signal"
	"github.com/d-sparks/gravy/signal/ipos"
	"github.com/d-sparks/gravy/signal/movingaverage"
	"github.com/d-sparks/gravy/signal/unlistings"
	"github.com/d-sparks/gravy/strategy"
	"github.com/d-sparks/gravy/strategy/buyandhold"
	"github.com/d-sparks/gravy/strategy/uniform"
)

// A TradingAlgorithm orchestrates stores, signals, strategies, and conducts trades.
type TradingAlgorithm struct {
	// For executing strategies and making trades.
	stores             map[string]db.Store
	signals            map[string]signal.Signal
	strategies         map[string]strategy.Strategy
	perStrategyMetrics map[string]metric.PerStrategyMetric
	exchange           gravy.Exchange

	// For debug.
	signalOrder            []string
	strategyOrder          []string
	perStrategyMetricOrder []string
	nonhiddenHeaders       stringset.StringSet
	headers                []string
	algorithmHeaders       []string
	debug                  map[string]string

	// Temporary. Store the successive product of the daily performance of each strategy.
	cumulativeReturn map[string]float64
}

func NewTradingAlgorithm(stores map[string]db.Store, exchange gravy.Exchange) TradingAlgorithm {
	// Initialize all members.
	algorithm := TradingAlgorithm{
		stores:             stores,
		signals:            map[string]signal.Signal{},
		strategies:         map[string]strategy.Strategy{},
		perStrategyMetrics: map[string]metric.PerStrategyMetric{},
		exchange:           exchange,

		nonhiddenHeaders: stringset.New(),
		signalOrder:      []string{},
		strategyOrder:    []string{},
		debug:            map[string]string{},

		cumulativeReturn: map[string]float64{},
	}

	// Initialize signals.
	algorithm.AddSignal(movingaverage.Name(100), movingaverage.New(100))
	algorithm.AddSignal(ipos.Name, ipos.New())
	algorithm.AddSignal(unlistings.Name, unlistings.New())

	// Initialize strategies.
	algorithm.AddStrategy(buyandhold.Name, buyandhold.New())
	algorithm.AddStrategy(uniform.Name, uniform.New())

	// Initialize metrics (must do these after all strategies).
	algorithm.AddPerStrategyMetric(ndayperf.Name(1), ndayperf.New(1))
	algorithm.AddPerStrategyMetric(ndayperf.Name(7), ndayperf.New(7))
	algorithm.AddPerStrategyMetric(ndayperf.Name(30), ndayperf.New(30))

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
	t.cumulativeReturn[name] = 1.0
}

// These must be added after all strategies are added.
func (t *TradingAlgorithm) AddPerStrategyMetric(name string, metric metric.PerStrategyMetric) {
	t.perStrategyMetrics[name] = metric
	t.perStrategyMetricOrder = append(t.perStrategyMetricOrder, name)
}

// Calculates data, signals, strategies, and executes trades.
func (t *TradingAlgorithm) Trade(date time.Time) error {
	// Clear debug output.
	t.debug = map[string]string{"date": date.Format("2006-01-02")}

	// Get current portfolio.
	// portfolio := t.exchange.CurrentPortfolio()

	// Get outputs of individual strategies.
	strategyOutputs := map[string]*strategy.StrategyOutput{}
	for name, strategy := range t.strategies {
		// Run strategy.
		strategyOutput, err := strategy.Run(date, t.stores, t.signals)
		if err != nil {
			return fmt.Errorf("Error evaluating strategy `%s`: `%s`", name, err.Error())
		}
		strategyOutputs[name] = strategyOutput

		// Calculate strategy metrics.
		for metricName, metric := range t.perStrategyMetrics {
			metricValue, err := metric.Value(date, t.stores, t.signals, strategyOutput)
			if err != nil {
				return fmt.Errorf("Error evaluating metric value `%s`: `%s`", metricName, err.Error())
			}
			t.debug[stratHeader(name, metricName)] = strconv.FormatFloat(metricValue, 'f', -1, 64)

			// This is a bit hacky, find a better place to do this.
			if metricName == ndayperf.Name(1) {
				t.cumulativeReturn[name] *= metricValue
				t.debug[stratHeader(name, "cumulative")] = strconv.FormatFloat(
					t.cumulativeReturn[name], 'f', -1, 64,
				)
			}
		}
	}

	// TODO Calculate orders
	return nil
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

		// Signal output headers.
		for _, signal := range t.signalOrder {
			for _, header := range t.signals[signal].Headers() {
				t.headers = append(t.headers, signalHeader(signal, header))
			}
		}

		// Strategy output headers.
		for _, strat := range t.strategyOrder {
			for _, header := range t.strategies[strat].Headers() {
				t.headers = append(t.headers, stratHeader(strat, header))
			}
			for _, metric := range t.perStrategyMetricOrder {
				t.headers = append(t.headers, algHeader(stratHeader(strat, metric)))
			}
			t.headers = append(t.headers, algHeader(stratHeader(strat, "cumulative")))
		}

		// Algorithm level output headers.
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
		if true || !hide || t.nonhiddenHeaders.Contains(algHeaderStr) {
			debug[algHeaderStr] = value
		}
	}

	return debug
}
