package movingaverage

// Rolling is a struct that tracks n-day averages for an asset or portfolio.
type Rolling struct {
	sum  float64
	n    int
	days int
}

// NewRolling creates a moving average struct that tracks the given windows.
func NewRolling(days int) *Rolling {
	return &Rolling{days: days}
}

// Observe inserts the value.
func (r *Rolling) Observe(x, x0 float64) {
	r.sum += x
	if r.n >= r.days {
		r.sum -= x0
	} else {
		r.n++
	}
}

// Value returns the value.
func (r *Rolling) Value() float64 {
	if r.n <= 0.0 {
		return 0.0
	}
	return r.sum / float64(r.n)
}
