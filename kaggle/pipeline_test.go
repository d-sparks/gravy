package kaggle

import (
	"testing"
	"time"

	"github.com/Clever/go-utils/stringset"
	"github.com/d-sparks/gravy/trading"
	"github.com/stretchr/testify/assert"
)

func TestInterpolateData(t *testing.T) {
	t1, err := time.Parse("2006-01-02", "2019-10-01")
	assert.NoError(t, err)
	t2, err := time.Parse("2006-01-02", "2019-10-02")
	assert.NoError(t, err)
	t3, err := time.Parse("2006-01-02", "2019-10-03")
	assert.NoError(t, err)

	dates := []time.Time{t1, t2, t3}

	window1 := trading.Window{
		Begin: t1,
		End:   t1,

		Close:  trading.Prices{"MSFT": 1.0},
		High:   trading.Prices{"MSFT": 2.0},
		Low:    trading.Prices{"MSFT": 3.0},
		Open:   trading.Prices{"MSFT": 4.0},
		Volume: trading.Prices{"MSFT": 5.0},

		Symbols: stringset.New("MSFT"),
	}
	window2 := trading.Window{
		Begin: t2,
		End:   t2,

		Close:  trading.Prices{},
		High:   trading.Prices{},
		Low:    trading.Prices{},
		Open:   trading.Prices{},
		Volume: trading.Prices{},

		Symbols: stringset.New(),
	}
	window3 := trading.Window{
		Begin: t3,
		End:   t3,

		Close:  trading.Prices{"MSFT": 2.0},
		High:   trading.Prices{"MSFT": 3.0},
		Low:    trading.Prices{"MSFT": 4.0},
		Open:   trading.Prices{"MSFT": 5.0},
		Volume: trading.Prices{"MSFT": 6.0},

		Symbols: stringset.New("MSFT"),
	}

	data := map[time.Time]*trading.Window{t1: &window1, t2: &window2, t3: &window3}

	InterpolateData(dates, data)
	assert.Equal(t, (1.0+2.0)/2.0, window2.Close["MSFT"])
	assert.Equal(t, (2.0+3.0)/2.0, window2.High["MSFT"])
	assert.Equal(t, (3.0+4.0)/2.0, window2.Low["MSFT"])
	assert.Equal(t, (4.0+5.0)/2.0, window2.Open["MSFT"])
	assert.Equal(t, (5.0+6.0)/2.0, window2.Volume["MSFT"])
}
