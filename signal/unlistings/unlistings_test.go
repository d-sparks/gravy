package unlistings

import (
	"testing"
	"time"

	"github.com/Clever/go-utils/stringset"
	"github.com/d-sparks/gravy/db"
	"github.com/d-sparks/gravy/db/dailyprices"
	"github.com/stretchr/testify/assert"
)

func TestUnlistings(t *testing.T) {
	todayTime, err := time.Parse("2006-01-02", "2013-10-28")
	assert.NoError(t, err)
	tomorrowTime, err := time.Parse("2006-01-02", "2013-10-29")
	assert.NoError(t, err)
	twoDaysTime, err := time.Parse("2006-01-02", "2013-10-29")
	assert.NoError(t, err)

	yesterday := db.Data{Tickers: stringset.New("GOOGL", "DELL")}
	today := db.Data{Tickers: stringset.New("GOOGL")}

	testStore := dailyprices.NewInMemoryStore()
	testStore.Set(todayTime, &yesterday)
	testStore.Set(tomorrowTime, &today)
	testStore.Set(twoDaysTime, &today)

	stores := map[string]db.Store{}
	stores[dailyprices.Name] = testStore

	unlistings := New()

	todayOutput, err := unlistings.Compute(todayTime, stores)
	assert.NoError(t, err)
	assert.True(t, todayOutput.StringSet.Equals(stringset.New("DELL")))

	tomorrowOutput, err := unlistings.Compute(tomorrowTime, stores)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(tomorrowOutput.StringSet))
}
