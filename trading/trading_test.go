package trading

import (
	"testing"

	"github.com/Clever/go-utils/stringset"
	"github.com/d-sparks/gravy/db"
	"github.com/stretchr/testify/assert"
)

func TestCapitalDistributionToAbstractPortfolioAndBack(t *testing.T) {
	tickers := stringset.New("MSFT", "GOOGL", "FB")
	prices := map[string]db.Prices{
		"MSFT":  db.Prices{Close: 3.0},
		"GOOGL": db.Prices{Close: 4.0},
		"FB":    db.Prices{Close: 2.5},
	}

	uniform := NewUniformCapitalDistribution(tickers)
	ap := uniform.ToAbstractPortfolioOnClose(prices, 1337.0)
	cd := ap.ToCapitalDistributionOnClose(prices)
	ap2 := cd.ToAbstractPortfolioOnClose(prices, 1337.0)

	assert.Equal(t, uniform, cd)
	assert.Equal(t, ap, ap2)
}
