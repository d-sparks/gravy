package supervisor

import (
	"context"
	"sync"
	"time"

	"github.com/d-sparks/gravy/registrar"
	supervisor_pb "github.com/d-sparks/gravy/supervisor/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// TradingMode represents the type of simulation or live trading supervisor is responsible for.
type TradingMode int

const (
	syncTM TradingMode = iota
	asyncTM
	paperTM
	liveTM
)

// S is the supervisor.
type S struct {
	supervisor_pb.UnimplementedSupervisorServer

	registrar *registrar.R

	// This is the source of truth. If trading live, S is in charge of keeping these up to date with the broker or
	// exchange.
	portfolios map[string]*supervisor_pb.Portfolio

	// Current trading mode.
	tradingMode TradingMode

	// For simulation
	tickTime time.Time

	// Locks and unlocks when simulations or trading modes are in progress.
	mu sync.Mutex
}

// PlaceOrder places an order in the set trading mode.
func (s *S) PlaceOrder(context.Context, *supervisor_pb.Order) (*supervisor_pb.OrderConfirmation, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PlaceOrder not implemented")
}

// GetPortfolio gets and returns the current portfolio for the requested algorithm.
func (s *S) GetPortfolio(context.Context, *supervisor_pb.AlgorithmId) (*supervisor_pb.Portfolio, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPortfolio not implemented")
}

// SynchronousDailySim kicks off a synchronous daily sim.
func (s *S) SynchronousDailySim(
	context.Context, *supervisor_pb.SynchronousDailySimInput,
) (*supervisor_pb.SynchronousDailySimOutput, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SynchronousDailySim not implemented")
}

// Abort aborts the current trading mode. For live trading, this should also try to intelligently close positions.
func (s *S) Abort(context.Context, *supervisor_pb.AbortInput) (*supervisor_pb.AbortOutput, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Abort not implemented")
}
