package dailyprices

import (
	"testing"
	"time"

	"github.com/d-sparks/gravy/db"
	"github.com/stretchr/testify/assert"
)

func TestInMemoryStoreValidDate(t *testing.T) {
	store := NewInMemoryStore()
	data := db.Data{}
	date, err := time.Parse("2006-01-02", "2020-01-01")
	assert.NoError(t, err)
	store.Set(date, &data)

	date2, err := time.Parse("2006-01-02", "2020-01-02")
	assert.NoError(t, err)

	dateValid, err := store.ValidDate(date)
	assert.NoError(t, err)
	assert.True(t, dateValid)

	date2Valid, err := store.ValidDate(date2)
	assert.NoError(t, err)
	assert.False(t, date2Valid)
}

func TestInMemoryStoreNextDate(t *testing.T) {
	store := NewInMemoryStore()
	data := db.Data{}

	date1, err := time.Parse("2006-01-02", "2020-01-01")
	assert.NoError(t, err)
	date2, err := time.Parse("2006-01-02", "2020-01-02")
	assert.NoError(t, err)
	date3, err := time.Parse("2006-01-02", "2020-01-03")
	assert.NoError(t, err)

	store.Set(date1, &data)
	store.Set(date3, &data)

	date1NextDate, err := store.NextDate(date1)
	assert.NoError(t, err)
	assert.Equal(t, date3, *date1NextDate)

	date2NextDate, err := store.NextDate(date2)
	assert.NoError(t, err)
	assert.Equal(t, date3, *date2NextDate)

	_, err = store.NextDate(date3)
	assert.Error(t, err)
}

func TestInMemoryStoreGet(t *testing.T) {
	store := NewInMemoryStore()
	data := db.Data{}

	date1, err := time.Parse("2006-01-02", "2020-01-01")
	assert.NoError(t, err)

	store.Set(date1, &data)

	retrieved, err := store.Get(date1)
	assert.NoError(t, err)
	assert.Equal(t, &data, retrieved)
}
