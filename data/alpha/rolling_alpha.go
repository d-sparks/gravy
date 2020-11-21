package alpha

import (
	"github.com/d-sparks/gravy/data/covariance"
	"github.com/d-sparks/gravy/data/variance"
)

// Rolling tracks alpha and beta over a rolling window.
type Rolling struct {
	cov   *covariance.Rolling
	varm  *variance.Rolling
	r     float64
	perf  float64
	bench float64
	n     int
	days  int
}

// NewRolling alpha / beta tracker for a prescribed number of days.
func NewRolling(days int, r float64) *Rolling {
	return &Rolling{cov: covariance.NewRolling(days), varm: variance.NewRolling(days), r: r, days: days}
}

// Observe observes a new value of the asset and benchmark prices.
func (r *Rolling) Observe(x, m, mux, mum, x0, m0, mux0, mum0 float64) {
	r.cov.Observe(x, m, mux, mum, x0, m0, mux0, mum0)
	r.varm.Observe(m, mum, m0, mum0)

	r.perf = (x - x0) / x0
	r.bench = (m - m0) / m0
	if r.n < r.days {
		r.n++
	}
}

// Beta returns the beta.
func (r *Rolling) Beta() float64 {
	return r.cov.Value() / r.varm.Value()
}

// Alpha returns the alpha.
func (r *Rolling) Alpha() float64 {
	return r.perf - r.r - r.Beta()*(r.bench-r.r)
}
