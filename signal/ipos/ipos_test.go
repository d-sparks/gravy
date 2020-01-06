package ipos

import (
	"testing"
	"time"

	"github.com/Clever/go-utils/stringset"
	"github.com/d-sparks/gravy/db"
	"github.com/d-sparks/gravy/db/dailyprices"
	"github.com/stretchr/testify/assert"
)

func TestIPOs(t *testing.T) {
	yesterdayTime, err := time.Parse("2006-01-02", "2004-08-18")
	assert.NoError(t, err)
	todayTime, err := time.Parse("2006-01-02", "2004-08-19")
	assert.NoError(t, err)

	yesterday := db.Data{Tickers: stringset.New("MSFT")}
	today := db.Data{Tickers: stringset.New("MSFT", "GOOGL")}

	testStore := dailyprices.NewInMemoryStore()
	testStore.Set(yesterdayTime, &yesterday)
	testStore.Set(todayTime, &today)

	stores := map[string]db.Store{}
	stores[dailyprices.Name] = testStore

	ipos := NewIPOs()

	yesterdayOutput, err := ipos.Compute(yesterdayTime, stores)
	assert.NoError(t, err)
	assert.True(t, yesterdayOutput.StringSet.Equals(stringset.New("MSFT")))

	todayOutput, err := ipos.Compute(todayTime, stores)
	assert.NoError(t, err)
	assert.True(t, todayOutput.StringSet.Equals(stringset.New("GOOGL")))
}
