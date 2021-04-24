package variance

import (
	"github.com/d-sparks/gravy/data/covariance"
	"github.com/d-sparks/gravy/data/mean"
)

// Rolling tracks the variance.
type Rolling struct {
	c *covariance.Rolling
}

// NewRolling creates a new variance.
func NewRolling(mux *mean.Rolling, days int) *Rolling {
	return &Rolling{c: covariance.NewRolling(mux, mux, days)}
}

// Observe observes a new value. Does not update the underlying moving averages, those should be updated separately.
func (v *Rolling) Observe(x float64) {
	v.c.Observe(x, x)
}

// Value returns the value of the sample variance.
func (v *Rolling) Value() float64 {
	return v.c.Value()
}

// UncorrectedValue returns the population variance.
func (v *Rolling) UncorrectedValue() float64 {
	return v.c.UncorrectedValue()
}
