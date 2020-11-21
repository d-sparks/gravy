package covariance

// Rolling captures a rolling covariance.
type Rolling struct {
	sumSqrs float64
	n       int
	days    int
}

// NewRolling makes a new rolling covariance for a prescribed number of days.
func NewRolling(days int) *Rolling {
	return &Rolling{days: days}
}

// Observe observes two values and updates the rolling covariance. Must pass in the values and averages of x, y that are
// falling off the back.
func (r *Rolling) Observe(x, y, mux, muy, x0, y0, mux0, muy0 float64) {
	r.sumSqrs += (x - mux) * (y - muy)
	if r.n >= r.days {
		r.sumSqrs -= (x0 - mux0) * (y0 - muy0)
	} else {
		r.n++
	}
}

// Value of the sample covariance.
func (r *Rolling) Value() float64 {
	if r.n <= 1 {
		return 0.0
	}
	return r.sumSqrs / (float64(r.n) - 1.0)
}

// ValueUncorrected is the sample covariance without Bessel's correction.
func (r *Rolling) UncorrectedValue() float64 {
	if r.n <= 0 {
		return 0.0
	}
	return r.sumSqrs / float64(r.n)
}
