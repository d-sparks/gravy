package ipos

import (
	"testing"
	"time"

	"github.com/Clever/go-utils/stringset"
	"github.com/d-sparks/gravy/db"
	"github.com/d-sparks/gravy/db/dailywindow"
	"github.com/d-sparks/gravy/trading"
	"github.com/stretchr/testify/assert"
)

func TestIPOs(t *testing.T) {
	yesterdayTime, err := time.Parse("2006-01-02", "2019-11-01")
	assert.NoError(t, err)
	todayTime, err := time.Parse("2006-01-02", "2019-11-02")
	assert.NoError(t, err)

	yesterday := trading.Window{Symbols: stringset.New("MSFT")}
	today := trading.Window{Symbols: stringset.New("MSFT", "GOOGL")}

	testStore := dailywindow.NewInMemoryStore()
	testStore.Set(yesterdayTime, &yesterday)
	testStore.Set(todayTime, &today)

	stores := map[string]db.Store{}
	stores[dailywindow.Name] = testStore

	IPOsSignal := New()

	yesterdayOutput := IPOsSignal.Compute(yesterdayTime, stores)
	assert.True(t, yesterdayOutput.StringSet.Equals(stringset.New("MSFT")))

	todayOutput := IPOsSignal.Compute(todayTime, stores)
	assert.True(t, todayOutput.StringSet.Equals(stringset.New("GOOGL")))
}
