package alpha

import (
	"fmt"
	"math"

	"github.com/d-sparks/gravy/data/covariance"
	"github.com/d-sparks/gravy/data/mean"
	"github.com/d-sparks/gravy/data/variance"
)

// Approximate number of trading days per year.
const daysPerYear float64 = 252.0

// An A tracks alpha and beta.
type A struct {
	x0 float64
	m0 float64
	r  float64
	n  float64

	rx *mean.M
	rm *mean.M

	cov  *covariance.C
	varm *variance.V
}

// New creates a new alpha/beta tracker.
func New(r float64) *A {
	return &A{r: r, rx: mean.New(), rm: mean.New(), cov: covariance.New(), varm: variance.New()}
}

// Observe observes the value and market value. Returns an error if either is 0.0 on the first observation. x is the
// value of the asset and m is the market/benchmark.
func (a *A) Observe(x float64, m float64) error {
	// Initialize if necessary.
	if a.n == 0.0 {
		if x <= 0.0 || m <= 0.0 {
			return fmt.Errorf("Cannot begin tracking alpha at 0.0")
		}
		a.x0 = x
		a.m0 = m
	}

	// Update statistics.
	a.n++
	a.cov.Observe(x, m)
	a.varm.Observe(m)
	a.rx.Observe(math.Pow(x/a.x0, daysPerYear/a.n))
	a.rm.Observe(math.Pow(m/a.m0, daysPerYear/a.n))

	return nil
}

// Beta returns the beta. Only valid after making an observation.
func (a *A) Beta() float64 {
	if a.n <= 0.0 {
		return 0.0
	}
	return a.cov.Value() / a.varm.Value()
}

// Alpha returns the alpha. Only valid after making an observation.
func (a *A) Alpha() float64 {
	if a.n <= 0.0 {
		return 0.0
	}
	return a.rx.Value() - a.r - a.Beta()*(a.rm.Value()-a.r)
}
