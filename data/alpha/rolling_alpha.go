package alpha

import (
	"github.com/d-sparks/gravy/data/covariance"
	"github.com/d-sparks/gravy/data/mean"
	"github.com/d-sparks/gravy/data/variance"
)

// Rolling tracks alpha and beta over a rolling window.
type Rolling struct {
	cov   *covariance.Rolling
	varm  *variance.Rolling
	r     float64
	perf  float64
	bench float64

	observations int // make sure we have two observations before returning a beta.
}

// NewRolling alpha / beta tracker for a prescribed number of days.
func NewRolling(mux *mean.Rolling, mum *mean.Rolling, days int, r float64) *Rolling {
	return &Rolling{
		cov:  covariance.NewRolling(mux, mum, days),
		varm: variance.NewRolling(mum, days),
		r:    r,
	}
}

// Observe observes a new value of the asset and benchmark prices.
func (r *Rolling) Observe(x float64, m float64) {
	r.cov.Observe(x, m)
	r.varm.Observe(m)

	r.perf = x
	r.bench = m

	r.observations++
}

// Beta returns the beta.
func (r *Rolling) Beta() float64 {
	if r.observations < 2 {
		return 0.0
	}
	return r.cov.Value() / r.varm.Value()
}

// Alpha returns the alpha.
func (r *Rolling) Alpha() float64 {
	return r.perf - r.r - r.Beta()*(r.bench-r.r)
}
