package variance

import "github.com/d-sparks/gravy/data/covariance"

// Streaming tracks the variance.
type Streaming struct {
	c *covariance.Streaming
}

// NewStreaming creates a new variance.
func NewStreaming() *Streaming {
	return &Streaming{c: covariance.NewStreaming()}
}

// Observe observes a new value. Returns error if the first value is 0.0.
func (s *Streaming) Observe(x float64) error {
	return s.c.Observe(x, x)
}

// Value returns the value of the variance.
func (s *Streaming) Value() float64 {
	return s.c.Value()
}

// UncorrectedValue returns the value without Bessel's correction.
func (s *Streaming) UncorrectedValue() float64 {
	return s.c.UncorrectedValue()
}

// RelativeValue returns the relative variance.
func (s *Streaming) RelativeValue() float64 {
	return s.c.RelativeValue()
}

// UncorrectedRelativeValue returns the relative variance without Bessel's correction.
func (s *Streaming) UncorrectedRelativeValue() float64 {
	return s.c.UncorrectedRelativeValue()
}
