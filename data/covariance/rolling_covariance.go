package covariance

import (
	"github.com/d-sparks/gravy/data/mean"
)

// Rolling tracks a rolling covariance. Can track multiple numbers of days.
type Rolling struct {
	cov *mean.Rolling
	mux *mean.Rolling
	muy *mean.Rolling

	n    float64
	days int
}

// NewRolling makes a new rolling covariance for a prescribed number of days.
func NewRolling(mux *mean.Rolling, muy *mean.Rolling, days int) *Rolling {
	return &Rolling{cov: mean.NewRolling(days), mux: mux, muy: muy, days: days}
}

// Observe observes two values and updates the rolling covariance. This does not update the underlying rolling
// averages and you should update them before updating the covariance.
func (r *Rolling) Observe(x float64, y float64) {
	r.cov.Observe((x - r.mux.Value(r.days)) * (y - r.muy.Value(r.days)))
	if r.n < float64(r.days) {
		r.n++
	}
}

// Value of the sample covariance (includes Bessel's correction).
func (r *Rolling) Value() float64 {
	if r.n <= 1.0 {
		return 0.0
	}
	return r.cov.Value(r.days) * r.n / (r.n - 1)
}
