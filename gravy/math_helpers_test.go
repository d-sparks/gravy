package gravy

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAnnualizedPerfHalfPeriod(t *testing.T) {
	init := 1.0
	mature := 1.1
	tradingDays := 252 / 2

	annualized := AnnualizedPerf(init, mature, tradingDays)
	expectedAnnualized := mature * mature
	assert.Equal(t, expectedAnnualized, annualized)
}

func TestAnnualizedPerfSinglePeriod(t *testing.T) {
	init := 1.0
	mature := 1.1
	tradingDays := 252

	annualized := AnnualizedPerf(init, mature, tradingDays)
	expectedAnnualized := mature
	assert.Equal(t, expectedAnnualized, annualized)
}

func TestAnnualizedPerfDoublePeriod(t *testing.T) {
	init := 1.0
	mature := 1.1
	tradingDays := 252 * 2

	annualized := AnnualizedPerf(init, mature, tradingDays)
	expectedAnnualized := math.Sqrt(mature)
	assert.Equal(t, expectedAnnualized, annualized)
}
