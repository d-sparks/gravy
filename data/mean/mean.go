package mean

// Streaming tracks the approximate mean of a stream of data.
type Streaming struct {
	mu float64
	n  float64
}

// NewStreaming returns an empty mean.
func NewStreaming() *Streaming {
	return &Streaming{}
}

// Observe observes a new value.
func (s *Streaming) Observe(x float64) {
	s.mu = ((s.n * s.mu) + x) / (s.n + 1.0)
	s.n += 1.0
}

// Value returns the mean.
func (s *Streaming) Value() float64 {
	return s.mu
}
