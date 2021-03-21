package covariance

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestCovariance tests the streaming covariance against a direct two-pass calculation.
func TestCovariance(t *testing.T) {
	covariance := NewStreaming()

	// Make two data sets.
	x := make([]float64, 10000)
	y := make([]float64, 10000)
	for i := 0; i < 10000; i++ {
		x[i] = 10000.0 + 5.0*float64(i)
		y[i] = x[i] + math.Sqrt(x[i])*math.Cos(float64(i))
	}

	// Calculate means directly.
	mux := 0.0
	muy := 0.0
	for i := range x {
		mux += x[i] / x[0]
		muy += y[i] / y[0]
	}
	mux /= float64(len(x))
	muy /= float64(len(y))

	// Calculate covariance both ways.
	cov := 0.0
	for i := 0; i < 10000; i++ {
		cov += (x[i]/x[0] - mux) * (y[i]/y[0] - muy)
		assert.Nil(t, covariance.Observe(x[i], y[i]))
	}
	cov /= 10000.0 - 1.0

	// Compare.
	assert.InDelta(t, cov, covariance.RelativeValue(), 1E-3)
}
