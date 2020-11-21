package variance

import "github.com/d-sparks/gravy/data/covariance"

// Rolling tracks the variance.
type Rolling struct {
	c *covariance.Rolling
}

// NewRolling creates a new variance.
func NewRolling(days int) *Rolling {
	return &Rolling{c: covariance.NewRolling(days)}
}

// Observe observes a new value. Returns error if the first value is 0.0.
func (v *Rolling) Observe(x, mux, x0, mux0 float64) {
	v.c.Observe(x, x, mux, mux, x0, x0, mux0, mux0)
}

// Value returns the value of the variance.
func (v *Rolling) Value() float64 {
	return v.c.Value()
}

// UncorrectedValue returns the value without Bessel's correction.
func (v *Rolling) UncorrectedValue() float64 {
	return v.c.UncorrectedValue()
}
