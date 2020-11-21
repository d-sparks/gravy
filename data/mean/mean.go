package mean

// M tracks the approximate mean of a stream of data.
type M struct {
	mu float64
	n  float64
}

// New returns an empty A.
func New() *M {
	return &M{}
}

// Observe observes a new value.
func (m *M) Observe(x float64) {
	m.mu = ((m.n * m.mu) + x) / (m.n + 1.0)
	m.n += 1.0
}

// Value returns the mean.
func (m *M) Value() float64 {
	return m.mu
}
