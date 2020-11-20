package covariance

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestCovariance tests the streaming covariance against a direct two-pass calculation.
func TestCovariance(t *testing.T) {
	covariance := New(10000.0, 5000.0)

	mux := 0.0
	muy := 0.0
	x := make([]float64, 10000)
	y := make([]float64, 10000)
	for i := 0; i < 10000; i++ {
		x[i] = 10000.0 + 5.0*float64(i)
		y[i] = x[i] + math.Sqrt(x[i])*math.Cos(float64(i))
		mux += x[i] / 10000.0
		muy += y[i] / 5000.0
		covariance.Observe(x[i], y[i])
	}
	mux /= 10000.0
	muy /= 10000.0

	cov := 0.0
	for i := 0; i < 10000; i++ {
		cov += (x[i]/10000.0 - mux) * (y[i]/5000.0 - muy)
	}
	cov /= 10000.0

	assert.InDelta(t, cov, covariance.RelativeValue(), 1E-2)
}
