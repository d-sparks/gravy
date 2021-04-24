package gravy

import (
	"math"
	"os"
	"sort"
	"time"

	"github.com/d-sparks/gravy/data/mean"
	"github.com/d-sparks/gravy/data/variance"
	"github.com/montanaflynn/stats"
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

// AnnualizedPerf extrapolates the annual return from an arbitrary number of trading days.
func AnnualizedPerf(init float64, mature float64, tradingDays int) float64 {
	return math.Pow(mature/init, 252.0/float64(tradingDays))
}

// Distribution attempts to describe a distribution.
type Distribution struct {
	Min         float64
	Max         float64
	Mean        float64
	StDev       float64
	Percentiles map[int]float64
}

// calculateDistributionImpl takes an iterable lambda and uses it to populate a distribution (sample or population).
func calculateDistributionImpl(
	begin int,
	end int,
	sample bool,
	lambda func(int) float64,
	percentiles ...int,
) *Distribution {
	output := Distribution{Percentiles: map[int]float64{}}

	// Use a rolling mean and variance, but recording all days.
	numVals := end - begin
	mu := mean.NewRolling(numVals)
	Var := variance.NewRolling(mu, numVals)
	for i := 0; i < numVals; i++ {
		mu.Observe(lambda(begin + i))
	}
	// Iterate separately on variance so that the fully calculated mean will be used.
	for i := 0; i < numVals; i++ {
		Var.Observe(lambda(begin + i))
	}

	// Fill output's mean and stdev.
	output.Mean = mu.Value(numVals)
	if sample {
		output.StDev = math.Sqrt(Var.Value())
	} else {
		output.StDev = math.Sqrt(Var.UncorrectedValue())
	}

	// Fill output's percentiles, min, and max.
	values := sort.Float64Slice(mu.GetBuffer())
	values.Sort()
	output.Min = values[0]
	output.Max = values[len(values)-1]
	for _, percentile := range percentiles {
		output.Percentiles[percentile], _ = stats.Percentile(stats.Float64Data(values), float64(percentile))
	}

	return &output
}

// CalculateDistribution fills a distribution where the iterable lambda is assumed to describe a population.
func CalculateDistribution(begin int, end int, lambda func(int) float64, percentiles ...int) *Distribution {
	return calculateDistributionImpl(begin, end, false, lambda, percentiles...)
}

// CalculateSampleDistribution fills a distribution where the iterable lambda is assumed to describe a sample.
func CalculateSampleDistribution(begin int, end int, lambda func(int) float64, percentiles ...int) *Distribution {
	return calculateDistributionImpl(begin, end, true, lambda, percentiles...)
}
