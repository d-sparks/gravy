package supervisor

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path"
	"sort"
	"strings"
	"sync"

	algorithmio_pb "github.com/d-sparks/gravy/algorithm/proto"
	dailyprices_pb "github.com/d-sparks/gravy/data/dailyprices/proto"
	"github.com/d-sparks/gravy/gravy"
	"github.com/d-sparks/gravy/registrar"
	supervisor_pb "github.com/d-sparks/gravy/supervisor/proto"
	"github.com/golang/protobuf/ptypes"
	timestamp_pb "github.com/golang/protobuf/ptypes/timestamp"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// TradingMode represents the type of simulation or live trading supervisor is responsible for.
type TradingMode int

const (
	// SyncTM represents Synchronous daily simulated trading.
	SyncTM TradingMode = iota
	// AsyncTM represents asynchronous simulated trading.
	AsyncTM
	// PaperTM represents live trading with paper money.
	PaperTM
	// LiveTM represents live trading with real money.
	LiveTM
)

// S is the supervisor.
type S struct {
	supervisor_pb.UnimplementedSupervisorServer

	// registrar for accessing algorithms and data.
	registrar *registrar.R

	// This is the source of truth. If trading live, S is in charge of keeping these up to date with the broker or
	// exchange.
	portfolios map[string]*supervisor_pb.Portfolio

	// Pending orders. For intraday, we may want orders to persist across ticks.
	pendingOrders []*supervisor_pb.Order

	// A channel for algorithms to signal that they are done.
	doneWorking chan struct{}

	// Current trading mode.
	tradingMode TradingMode

	// Logging.
	algorithmsOut      map[string]*bufio.Writer
	algorithmsOutOrder []string
	out                *bufio.Writer

	// Locks and unlocks when simulations or trading modes are in progress.
	mu sync.Mutex
}

// New creates a new supervisor in the given trading mode.
func New(tradingMode TradingMode) (*S, error) {
	var s S

	s.doneWorking = make(chan struct{})
	s.tradingMode = tradingMode

	return &s, nil
}

// Init initializes the supervisor, in particular the registrar.
func (s *S) Init() error {
	var err error
	s.registrar, err = registrar.New()
	if err != nil {
		return fmt.Errorf("Error constructing registrar: %s", err.Error())
	}
	return nil
}

// Close closes the registrar.
func (s *S) Close() {
	s.registrar.Close()
}

// PlaceOrder places an order in the set trading mode.
func (s *S) PlaceOrder(ctx context.Context, order *supervisor_pb.Order) (*supervisor_pb.OrderConfirmation, error) {
	s.pendingOrders = append(s.pendingOrders, order)
	return &supervisor_pb.OrderConfirmation{}, nil
}

// GetPortfolio gets and returns the current portfolio for the requested algorithm.
func (s *S) GetPortfolio(
	ctx context.Context,
	algorithmID *supervisor_pb.AlgorithmId,
) (*supervisor_pb.Portfolio, error) {
	portfolio, ok := s.portfolios[algorithmID.GetAlgorithmId()]
	if !ok {
		return nil, fmt.Errorf("No portfolio for algorith: `%s`", algorithmID.GetAlgorithmId())
	}
	return portfolio, nil
}

// initPortfolios makes each portfolio have no positions and an initial capital in USD.
func (s *S) initPortfolios(capitalEach float64) {
	s.portfolios = map[string]*supervisor_pb.Portfolio{}
	for algorithmID := range s.registrar.Algorithms {
		s.portfolios[algorithmID] = &supervisor_pb.Portfolio{}
		s.portfolios[algorithmID].Stocks = map[string]float64{}
		s.portfolios[algorithmID].Usd = capitalEach
	}
}

// getTradingDatesInRange returns the trading dates in the given range.
func (s *S) getTradingDatesInRange(
	ctx context.Context, input *supervisor_pb.SynchronousDailySimInput,
) (*dailyprices_pb.TradingDates, error) {
	var simRange dailyprices_pb.Range
	simRange.Lb = input.GetStart()
	simRange.Ub = input.GetEnd()
	return s.registrar.DailyPrices.TradingDatesInRange(ctx, &simRange)
}

// executeAllAlgorithms tells each algorithm to begin working.
func (s *S) executeAllAlgorithms(ctx context.Context, algorithmDate *timestamp_pb.Timestamp) error {
	var input algorithmio_pb.Input
	input.Timestamp = algorithmDate
	for _, algorithm := range s.registrar.Algorithms {
		go algorithm.Execute(ctx, &input)
	}
	return nil
}

// logOrder logs an order to the algorithm's output file.
func (s *S) logOrder(algorithmID string, filled bool, buySell string, volume float64, price float64, ticker string) {
	if !filled {
		buySell = "(FAILED ORDER) " + buySell
	}
	buySell += fmt.Sprintf(" %f units of %s at %f [value: %f]\n", volume, ticker, price, volume*price)
	s.algorithmsOut[algorithmID].WriteString(buySell)
}

// handleOrder attempts to handle order and returns true if the order is fulfilled. Does not yet support short selling.
func (s *S) handleOrder(order *supervisor_pb.Order, dailyPrices *dailyprices_pb.DailyPrices) bool {
	// Check if the portfolio has enough of the stock to sell or enough USD to cover the limit price of the order.
	algorithmID := order.GetAlgorithmId().GetAlgorithmId()
	portfolio := s.portfolios[algorithmID]
	ticker := order.GetTicker()
	amountHeld, holding := portfolio.GetStocks()[ticker]
	if order.GetVolume() < 0 && (!holding || amountHeld < order.GetVolume()) {
		// Don't have enough of the stock to sell.
		return false
	} else if order.GetVolume()*order.GetLimit() > portfolio.GetUsd() {
		// Don't have enough money to pay for the stocks.
		return false
	}

	prices, ok := dailyPrices.GetStockPrices()[ticker]
	if !ok {
		// This stock isn't on the market.
		return false
	}
	open, close := prices.GetOpen(), prices.GetClose()
	price := (open + close) / 2.0

	if order.GetVolume() > 0 && price <= order.GetLimit() {
		// We checked above that volume * limit <= $USD, and since here we have conditioned on price <= limit,
		// we know volume * price <= $USD as well.
		portfolio.Usd -= order.GetVolume() * price
		portfolio.Stocks[ticker] += order.GetVolume()
		s.logOrder(algorithmID, true, "BUY", order.GetVolume(), price, ticker)
		return true
	} else if order.GetVolume() < 0 && price >= order.GetStop() {
		// We checked already that we have enough stock to sell.
		portfolio.Usd += order.GetVolume() * price
		portfolio.Stocks[ticker] -= order.GetVolume()
		s.logOrder(algorithmID, true, "SELL", order.GetVolume(), price, ticker)
		return true
	}

	// TODO: Handle the case of non-expiring orders.
	return false
}

// fillPendingOrders attempts to fill all pending orders.
func (s *S) fillPendingOrders(
	ctx context.Context,
	tradingDate *timestamp_pb.Timestamp,
) (*dailyprices_pb.DailyPrices, error) {
	var pricesReq dailyprices_pb.Request
	pricesReq.Timestamp = tradingDate
	pricesReq.Version = 0
	tradingPrices, err := s.registrar.DailyPrices.Get(ctx, &pricesReq)
	if err != nil {
		return nil, fmt.Errorf("Error getting trading prices: %s", err.Error())
	}
	for _, order := range s.pendingOrders {
		// TODO: Check the return value of this to see whether to continue the order into the next tick.
		s.handleOrder(order, tradingPrices)
	}
	s.pendingOrders = []*supervisor_pb.Order{}
	return tradingPrices, nil
}

// closeDelistedPositions compares the stocks listed in two DailyPrices. For stocks in the former but not in the latter,
// look for any portfolios that hold that stock and close the position at the former DailyPrices closing price.
func (s *S) closeDelistedPositions(
	ctx context.Context,
	algorithmDate *timestamp_pb.Timestamp,
	tradingPrices *dailyprices_pb.DailyPrices,
) error {
	// Get prices as the algorithm is aware.
	var pricesReq dailyprices_pb.Request
	pricesReq.Timestamp = algorithmDate
	pricesReq.Version = 0
	algorithmPrices, err := s.registrar.DailyPrices.Get(ctx, &pricesReq)
	if err != nil {
		return fmt.Errorf("Error getting algorithm trading prices: %s", err.Error())
	}

	// Find stocks the algorithm was holding but aren't in the tradingPrices.
	delistings := map[string]float64{}
	for ticker := range algorithmPrices.GetStockPrices() {
		_, ok := tradingPrices.GetStockPrices()[ticker]
		if !ok {
			delistings[ticker] = algorithmPrices.GetStockPrices()[ticker].GetClose()
		}
	}

	// For each delisted stock and each portfolio that holds it, close the position.
	for ticker, closingPrice := range delistings {
		for algorithmID, portfolio := range s.portfolios {
			volume, ok := portfolio.GetStocks()[ticker]
			if ok && volume > 0 {
				portfolio.Stocks[ticker] = 0
				portfolio.Usd += volume * closingPrice
				s.algorithmsOut[algorithmID].WriteString(
					fmt.Sprintf(
						"%s delisted, closed %f units at %f [value: %f]\n",
						ticker,
						volume,
						closingPrice,
						volume*closingPrice,
					),
				)
			}
		}
	}

	return nil
}

// setupOutput makes output files for each algorithm and for gravy and returns them in a map. Also creates the
// bufio.Writers used in the class.
func (s *S) setupOutput(dir string) (closer func(), err error) {
	// Create output directory if it doesn't exist.
	_, err = os.Stat(dir)
	if os.IsNotExist(err) {
		if err = os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("Error creating directory `%s`: %s", dir, err.Error())
		}
	}

	// Gravy log.
	files := map[string]*os.File{}
	files["gravy"], err = os.Create(path.Join(dir, "gravy.log"))
	if err != nil {
		return nil, fmt.Errorf("Error creating gravy log: %s", err.Error())
	}
	s.out = bufio.NewWriter(files["gravy"])

	// Create log for each algorithm.
	s.algorithmsOut = map[string]*bufio.Writer{}
	for algorithmID := range s.registrar.Algorithms {
		files[algorithmID], err = os.Create(path.Join(dir, algorithmID+".csv"))
		if err != nil {
			return nil, fmt.Errorf("Error creating log for algo `%s`: %s", algorithmID, err.Error())
		}
		s.algorithmsOut[algorithmID] = bufio.NewWriter(files[algorithmID])
		s.algorithmsOutOrder = append(s.algorithmsOutOrder, algorithmID)
	}
	sort.Strings(s.algorithmsOutOrder)

	// Write headers for gravy log.
	s.out.WriteString(strings.Join(append([]string{"date"}, s.algorithmsOutOrder...), ",") + "\n")

	// Create closer
	closer = func() {
		s.out.Flush()
		files["gravy"].Close()
		for algorithmID := range s.registrar.Algorithms {
			s.algorithmsOut[algorithmID].Flush()
			files[algorithmID].Close()
		}
	}

	return
}

// logTick logs data for a tick to the gravy log.
func (s *S) logTick(timestamp *timestamp_pb.Timestamp, prices *dailyprices_pb.DailyPrices) error {
	// Convert to native timestamp.
	nativeTime, err := ptypes.Timestamp(timestamp)
	if err != nil {
		return fmt.Errorf("Error converting timestamp to native time type: %s", err.Error())
	}

	// Get columns.
	var cols []string = make([]string, len(s.algorithmsOutOrder)+1)
	cols[0] = nativeTime.Format("2006-01-02")
	for i, algorithmID := range s.algorithmsOutOrder {
		cols[i+1] = fmt.Sprintf("%f", gravy.PortfolioValue(s.portfolios[algorithmID], prices))
	}

	// Log
	s.out.WriteString(strings.Join(cols, ",") + "\n")

	return nil
}

// SynchronousDailySim kicks off a synchronous daily sim.
func (s *S) SynchronousDailySim(
	ctx context.Context,
	input *supervisor_pb.SynchronousDailySimInput,
) (*supervisor_pb.SynchronousDailySimOutput, error) {
	// Lock the mutex.
	s.mu.Lock()
	defer s.mu.Unlock()

	// Initialize algorithms and establish a liquid portfolio for each algorithm. TODO: parametrize the capital.
	s.registrar.InitAllAlgorithms()
	defer s.registrar.CloseAllAlgorithms()
	s.initPortfolios(1000 * 1000 * 1000)

	// Create output files.
	closer, err := s.setupOutput(input.GetOutputDir())
	if err != nil {
		return nil, fmt.Errorf("Error creating log file: %s", err.Error())
	}
	defer closer()

	// Get trading dates in range.
	tradingDates, err := s.getTradingDatesInRange(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("Error getting trading dates: %s", err.Error())
	}

	// Simulate over the trading dates.
	for i := 1; i < len(tradingDates.GetTimestamps()); i++ {
		// The algorithms will have the prices from the previous tick, but trades are executed against prices
		// in the current tick.
		algorithmDate := tradingDates.GetTimestamps()[i-1]
		tradingDate := tradingDates.GetTimestamps()[i]

		// Tell each algorithm to begin working.
		if err := s.executeAllAlgorithms(ctx, algorithmDate); err != nil {
			return nil, err
		}

		// Wait for each algorithm to be done, up to a max timeout. TODO: handle deadlock gracefully. For
		// example have another goroutine that periodically checks if we're timed out, and if so, closes the
		// channel.
		for range s.registrar.Algorithms {
			<-s.doneWorking
		}

		// Try to fulfill pending orders.
		tradingPrices, err := s.fillPendingOrders(ctx, tradingDate)
		if err != nil {
			return nil, err
		}

		// Check if any stocks have been delisted. If they have, close out the position for algorithms holding
		// that stock.
		s.closeDelistedPositions(ctx, algorithmDate, tradingPrices)

		// Log for this tick.
		s.logTick(tradingDate, tradingPrices)
	}

	var synchronousDailySimOutput supervisor_pb.SynchronousDailySimOutput
	return &synchronousDailySimOutput, nil
}

// DoneTrading lets an algorithm signal that it is done trading in this tick.
func (s *S) DoneTrading(
	ctx context.Context,
	algorithmID *supervisor_pb.AlgorithmId,
) (*supervisor_pb.DoneTradingResponse, error) {
	s.doneWorking <- struct{}{}
	return &supervisor_pb.DoneTradingResponse{}, nil
}

// Abort aborts the current trading mode. For live trading, this should also try to intelligently close positions.
func (s *S) Abort(context.Context, *supervisor_pb.AbortInput) (*supervisor_pb.AbortOutput, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Abort not implemented")
}