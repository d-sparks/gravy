package covariance

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
type C struct {
	meanxnorm float64
	meanynorm float64
	n         float64
	c         float64

	normx float64
	normy float64
}

// New creates a new covariance from an initial observation. Panics if these numbers aren't positive.
func New(x float64, y float64) *C {
	if x <= 0.0 || y <= 0.0 {
		panic("Can't track relative covariance starting at 0.0")
	}
	c := C{normx: x, normy: y}
	c.Observe(x, y)
	return &c
}

// Observe observes two values if the observation hasn't been recorded at this time.
func (c *C) Observe(x float64, y float64) {
	// Update covariance.
	c.n++
	dx := ((x / c.normx) - c.meanxnorm)
	dy := ((y / c.normy) - c.meanynorm)
	c.meanxnorm += dx / c.n
	c.meanynorm += dy / c.n
	c.c += dx * dy
}

// Value returns the current covariance.
func (c *C) Value() float64 {
	if c.n <= 1 {
		return c.c * c.normx * c.normy
	}
	return c.c * c.normx * c.normy / (c.n - 1.0)
}

// UncorrectedValue returns the value without Bessel's correction.
func (c *C) UncorrectedValue() float64 {
	return c.c * c.normx * c.normy / c.n
}

// RelativeValue returns the relative covariance.
func (c *C) RelativeValue() float64 {
	if c.n <= 1 {
		return c.c
	}
	return c.c / (c.n - 1.0)
}

// UncorrectedRelativeValue returns the relative covariance without Bessel's correction.
func (c *C) UncorrectedRelativeValue() float64 {
	return c.c / c.n
}
