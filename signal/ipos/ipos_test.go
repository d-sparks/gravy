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
	todayTime, err := time.Parse("2006-01-02", "2004-08-18")
	assert.NoError(t, err)
	tomorrowTime, err := time.Parse("2006-01-02", "2004-08-19")
	assert.NoError(t, err)
	twoDaysTime, err := time.Parse("2006-01-02", "2004-08-19")
	assert.NoError(t, err)

	today := db.Data{Tickers: stringset.New("MSFT")}
	tomorrow := db.Data{Tickers: stringset.New("MSFT", "GOOGL")}

	testStore := dailyprices.NewInMemoryStore()
	testStore.Set(todayTime, &today)
	testStore.Set(tomorrowTime, &tomorrow)
	testStore.Set(twoDaysTime, &tomorrow)

	stores := map[string]db.Store{}
	stores[dailyprices.Name] = testStore

	ipos := New()

	todayOutput, err := ipos.Compute(todayTime, stores)
	assert.NoError(t, err)
	assert.True(t, todayOutput.StringSet.Equals(stringset.New("GOOGL")))

	tomorrowOutput, err := ipos.Compute(tomorrowTime, stores)
	assert.NoError(t, err)
	assert.True(t, tomorrowOutput.StringSet.Equals(stringset.New()))
}
