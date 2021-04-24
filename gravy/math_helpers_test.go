package gravy

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestAnnualizedPerfHalfPeriod tests that a year's return is the square of 6 months' return.
func TestAnnualizedPerfHalfPeriod(t *testing.T) {
	init := 1.0
	mature := 1.1
	tradingDays := 252 / 2

	annualized := AnnualizedPerf(init, mature, tradingDays)
	expectedAnnualized := mature * mature
	assert.Equal(t, expectedAnnualized, annualized)
}

// TestAnnualizedPerfSinglePeriod tests that a year's return is year's return.
func TestAnnualizedPerfSinglePeriod(t *testing.T) {
	init := 1.0
	mature := 1.1
	tradingDays := 252

	annualized := AnnualizedPerf(init, mature, tradingDays)
	expectedAnnualized := mature
	assert.Equal(t, expectedAnnualized, annualized)
}

// TestAnnualizedPerfDoublePeriod tests that a year's return is the sqrt of two years' return.
func TestAnnualizedPerfDoublePeriod(t *testing.T) {
	init := 1.0
	mature := 1.1
	tradingDays := 252 * 2

	annualized := AnnualizedPerf(init, mature, tradingDays)
	expectedAnnualized := math.Sqrt(mature)
	assert.Equal(t, expectedAnnualized, annualized)
}

// TestCalculateDistribution roughly tests the various statistics in calculate distribution.
func TestCalculateDistribution(t *testing.T) {
	data := []float64{1.0, 3.0, 7.0, 9.0, 10.0, 2.0, 8.0, 4.0, 6.0, 5.0}
	lambda := func(i int) float64 {
		return data[i]
	}

	distribution := CalculateDistribution(0, len(data), lambda, 20, 50, 70)
	assert.Equal(t, 1.0, distribution.Min, "Incorrect min")
	assert.Equal(t, 10.0, distribution.Max, "Incorrect max")
	assert.Equal(t, 5.5, distribution.Mean, "Incorrect mean")
	assert.Less(t, math.Abs(distribution.StDev-2.8722813232690143), 1e-6, "Incorrect stdev")
	assert.Equal(t, 2.0, distribution.Percentiles[20])
	assert.Equal(t, 5.0, distribution.Percentiles[50])
	assert.Equal(t, 7.0, distribution.Percentiles[70])
}
