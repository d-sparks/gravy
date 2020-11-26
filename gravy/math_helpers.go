package gravy

import "math"

// DivideOrZero returns the ratio of the two numbers unless the denominator is very small, in which case returns 0.0.
func DivideOrZero(n float64, d float64) float64 {
	if math.Abs(d) < 1E-6 {
		return 0.0
	}
	return n / d
}

// RelativePerfOrZero returns the relative performance or zero if the initial value is too small.
func RelativePerfOrZero(p float64, p0 float64) float64 {
	return DivideOrZero(p-p0, p0)
}
