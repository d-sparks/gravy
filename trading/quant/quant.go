package quant

import "math"

type StatStream struct {
	n  float64
	mu float64
	m2 float64
}

// https://en.wikipedia.org/wiki/Algorithms_for_calculating_variance#Welford's_online_algorithm
func NewStatStream() StatStream {
	return StatStream{
		n:  0.0,
		mu: 0.0,
		m2: 0.0, // \sum_{i=1...n} (x_i - mu)^2
	}
}

func (s *StatStream) Mu() float64 {
	return s.mu
}

func (s *StatStream) Sigma() float64 {
	return math.Sqrt(s.m2 / (s.n - 1.0))
}

func (s *StatStream) Observe(val float64) {
	muPrev := s.mu
	s.mu = (s.n*muPrev + val) / (s.n + 1.0)
	s.m2 += (val - muPrev) * (val - s.mu)
	s.n += 1.0
}
