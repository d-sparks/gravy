package covariance

import "fmt"

// Covariance implements this streaming algorithm for covariance I found on Wikipedia:
//
// def online_covariance(data1, data2):
//     meanx = meany = C = n = 0
//     for x, y in zip(data1, data2):
//         n += 1
//         dx = x - meanx
//         meanx += dx / n
//         meany += (y - meany) / n
//         C += dx * (y - meany)
//
//     population_covar = C / n
//     # Bessel's correction for sample variance
//     sample_covar = C / (n - 1)
//
type Streaming struct {
	meanxnorm float64
	meanynorm float64
	n         float64
	c         float64

	normx float64
	normy float64
}

// NewStreaming creates a new covariance.
func NewStreaming() *Streaming {
	return &Streaming{}
}

// Observe observes two values if the observation hasn't been recorded at this time.
func (s *Streaming) Observe(x float64, y float64) error {
	// Record first valid observation.
	if s.n == 0.0 {
		if x <= 0.0 || y <= 0.0 {
			return fmt.Errorf("Cannot start tracking covariance from 0.0")
		}
		s.normx = x
		s.normy = y
	}

	// Update covariance.
	s.n++
	dx := ((x / s.normx) - s.meanxnorm)
	dy := ((y / s.normy) - s.meanynorm)
	s.meanxnorm += dx / s.n
	s.meanynorm += dy / s.n
	s.c += dx * dy

	return nil
}

// Value returns the current covariance.
func (s *Streaming) Value() float64 {
	if s.n <= 1 {
		return s.c * s.normx * s.normy
	}
	return s.c * s.normx * s.normy / (s.n - 1.0)
}

// UncorrectedValue returns the value without Bessel's correction.
func (s *Streaming) UncorrectedValue() float64 {
	return s.c * s.normx * s.normy / s.n
}

// RelativeValue returns the relative covariance.
func (s *Streaming) RelativeValue() float64 {
	if s.n <= 1 {
		return s.c
	}
	return s.c / (s.n - 1.0)
}

// UncorrectedRelativeValue returns the relative covariance without Bessel's correction.
func (s *Streaming) UncorrectedRelativeValue() float64 {
	if s.n <= 0 {
		return s.c
	}
	return s.c / s.n
}
