package dailypricespipeline

import (
	"testing"

	"github.com/Clever/go-utils/stringset"
	"github.com/stretchr/testify/assert"
)

func TestInterpolateData(t *testing.T) {
	const (
		jan1  = "2019-01-01"
		jan2  = "2019-01-02"
		jan3  = "2019-01-03"
		jan4  = "2019-01-04"
		MSFT  = "MSFT"
		GOOGL = "GOOGL"
	)
	dates := []string{jan1, jan2, jan3}
	dateToTickerToCols := map[string]map[string][]float64{
		jan1: map[string][]float64{
			MSFT:  []float64{1.0, 10.0},
			GOOGL: []float64{2.0, 8.0},
		},
		jan2: map[string][]float64{
			GOOGL: []float64{3.0, 16.0},
		},
		jan3: map[string][]float64{
			MSFT:  []float64{3.0, 0.0},
			GOOGL: []float64{2.0, 8.0},
		},
		jan4: map[string][]float64{
			GOOGL: []float64{2.0, 8.0},
		},
	}
	dateToTickers := map[string]stringset.StringSet{
		jan1: stringset.New(MSFT, GOOGL),
		jan2: stringset.New(GOOGL),
		jan3: stringset.New(MSFT, GOOGL),
		jan4: stringset.New(GOOGL),
	}

	InterpolateData(dates, dateToTickerToCols, dateToTickers)

	// Expect interpolation for MSFT.
	assert.Equal(t, (1.0+3.0)/2.0, dateToTickerToCols[jan2][MSFT][0])
	assert.Equal(t, (10.0+0.0)/2.0, dateToTickerToCols[jan2][MSFT][1])
	assert.True(t, dateToTickers[jan2].Contains(MSFT))

	// Expect no interpolateion for GOOGL.
	assert.Equal(t, 3.0, dateToTickerToCols[jan2][GOOGL][0])
	assert.Equal(t, 16.0, dateToTickerToCols[jan2][GOOGL][1])

	// Expect no interpolation for MSFT after its last listing.
	_, ok := dateToTickerToCols[jan4][MSFT]
	assert.False(t, ok)
}