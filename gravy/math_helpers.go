package gravy

import (
	"math"
	"os"
	"time"
)

// DivideOrZero returns the ratio of the two numbers unless the denominator is very small, in which case returns 0.0.
func DivideOrZero(n, d float64) float64 {
	if math.Abs(d) < 1E-6 {
		return 0.0
	}
	return n / d
}

// RelativePerfOrZero returns the relative performance or zero if the initial value is too small.
func RelativePerfOrZero(p, p0 float64) float64 {
	return DivideOrZero(p-p0, p0)
}

// ZScore is a convenience method for zscore. Returns 0 if standard deviation is very close to 0.0.
func ZScore(x, mu, sigma float64) float64 {
	return DivideOrZero(x-mu, sigma)
}

// TimePIDSeed is a seed for random numbers: sensitive to process ID and current unix timestamp.
func TimePIDSeed() int64 {
	now := time.Now().Unix()
	pid := int64(os.Getpid())
	return now % (pid * pid)
}
