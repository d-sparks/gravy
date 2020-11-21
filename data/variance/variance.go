package variance

import "github.com/d-sparks/gravy/data/covariance"

// V tracks the variance.
type V struct {
	c *covariance.C
}

// New creates a new variance.
func New() *V {
	return &{c: covariance.New()}
}

// Observe observes a new value. Returns error if the first value is 0.0.
func (v *V) Observe(x float64) error {
	return v.c.Observe(x, x)
}

// Value returns the value of the variance.
func (v *V) Value() float64 {
	return v.c.Value()
}

