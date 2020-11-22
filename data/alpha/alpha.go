package alpha

import (
	"fmt"
	"math"

	"github.com/d-sparks/gravy/data/covariance"
	"github.com/d-sparks/gravy/data/mean"
	"github.com/d-sparks/gravy/data/variance"
)

// Approximate number of trading days per year.
const daysPerYear float64 = 252.0

// Streaming tracks alpha and beta for a stream of data. (Relatively low memory overhead.)
type Streaming struct {
	x0 float64
	m0 float64
	r  float64
	n  float64

	rx *mean.Streaming
	rm *mean.Streaming

	cov  *covariance.Streaming
	varm *variance.Streaming
}

// New creates a new alpha/beta tracker.
func NewStreaming(r float64) *Streaming {
	return &Streaming{
		r:    r,
		rx:   mean.NewStreaming(),
		rm:   mean.NewStreaming(),
		cov:  covariance.NewStreaming(),
		varm: variance.NewStreaming(),
	}
}

// Observe observes the value and market value. Returns an error if either is 0.0 on the first observation. x is the
// value of the asset and m is the market/benchmark.
func (s *Streaming) Observe(x float64, m float64) error {
	// Initialize if necessary.
	if s.n == 0.0 {
		if x <= 0.0 || m <= 0.0 {
			return fmt.Errorf("Cannot begin tracking alpha at 0.0")
		}
		s.x0 = x
		s.m0 = m
	}

	// Update statistics.
	s.n++
	s.cov.Observe(x, m)
	s.varm.Observe(m)
	s.rx.Observe(math.Pow(x/s.x0, daysPerYear/s.n))
	s.rm.Observe(math.Pow(m/s.m0, daysPerYear/s.n))

	return nil
}

// Beta returns the beta. Only valid after making an observation.
func (s *Streaming) Beta() float64 {
	if s.n <= 0.0 {
		return 0.0
	}
	return s.cov.Value() / s.varm.Value()
}

// Alpha returns the alphs. Only valid after making an observation.
func (s *Streaming) Alpha() float64 {
	if s.n <= 0.0 {
		return 0.0
	}
	return s.rx.Value() - s.r - s.Beta()*(s.rm.Value()-s.r)
}
