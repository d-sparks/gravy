package covariance

import (
	"testing"

	"github.com/d-sparks/gravy/data/movingaverage"
	"github.com/stretchr/testify/assert"
)

// TestRollingCovariance tests the rolling covariance implementation against direct calculations.
func TestRollingCovariance(t *testing.T) {
	// Instantiate classes.
	mux := movingaverage.New(20)
	muy := movingaverage.New(20)
	cov := NewRolling(mux, muy, 20)

	// Make some data.
	x := make([]float64, 40)
	y := make([]float64, 40)
	for i := 0; i < 40; i++ {
		x[i] = float64(i)
		y[i] = 1.5 * float64(i%2) * float64(i)
	}

	// Observe with implementation.
	for i := 0; i < 40; i++ {
		mux.Observe(x[i])
		muy.Observe(y[i])
		cov.Observe(x[i], y[i])
	}

	// Calculate directly
	mux = movingaverage.New(40)
	muy = movingaverage.New(40)
	sumSqrs := 0.0
	for i := 0; i < 40; i++ {
		mux.Observe(x[i])
		muy.Observe(y[i])
		if i >= 20 {
			sumSqrs += (x[i] - mux.Value(20)) * (y[i] - muy.Value(20))
		}
	}

	assert.Equal(t, sumSqrs/(20.0-1.0), cov.Value())

}

// TestRollingCovarianceUnderfull tests the rolling covariance implementation against direct calculations.
func TestRollingCovarianceUnderfull(t *testing.T) {
	// Instantiate classes.
	mux := movingaverage.New(20)
	muy := movingaverage.New(20)
	cov := NewRolling(mux, muy, 20)

	// Make some data.
	x := make([]float64, 5)
	y := make([]float64, 5)
	for i := 0; i < 5; i++ {
		x[i] = float64(i)
		y[i] = 1.5 * float64(i%2) * float64(i)
	}

	// Observe with implementation.
	for i := 0; i < 5; i++ {
		mux.Observe(x[i])
		muy.Observe(y[i])
		cov.Observe(x[i], y[i])
	}

	// Calculate directly
	mux = movingaverage.New(40)
	muy = movingaverage.New(40)
	sumSqrs := 0.0
	for i := 0; i < 5; i++ {
		mux.Observe(x[i])
		muy.Observe(y[i])
		sumSqrs += (x[i] - mux.Value(20)) * (y[i] - muy.Value(20))
	}

	assert.Equal(t, sumSqrs/(5.0-1.0), cov.Value())

}
