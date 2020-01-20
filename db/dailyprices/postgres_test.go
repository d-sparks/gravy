package dailyprices

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	dbURL       = "postgres://localhost/gravy?sslmode=disable"
	pricesTable = "dailyprices"
	datesTable  = "tradingdates"
)

func getDB(t *testing.T) *PostgresStore {
	store, err := NewPostgresStore(dbURL, pricesTable, datesTable)
	assert.NoError(t, err)
	return store
}

func TestPostgresStoreValidDate(t *testing.T) {
	store := getDB(t)
	defer store.Close()

	jan3, err := time.Parse("2006-01-02", "2010-01-03")
	assert.NoError(t, err)
	jan4, err := time.Parse("2006-01-02", "2010-01-04")
	assert.NoError(t, err)

	jan3Valid, err := store.ValidDate(jan3)
	assert.NoError(t, err)
	assert.False(t, jan3Valid)

	jan4Valid, err := store.ValidDate(jan4)
	assert.NoError(t, err)
	assert.True(t, jan4Valid)
}

func TestPostgresStoreNextDate(t *testing.T) {
	store := getDB(t)
	defer store.Close()

	dec31, err := time.Parse("2006-01-02", "2009-12-31")
	assert.NoError(t, err)
	jan3, err := time.Parse("2006-01-02", "2010-01-03")
	assert.NoError(t, err)
	jan4, err := time.Parse("2006-01-02", "2010-01-04")
	assert.NoError(t, err)
	jan5, err := time.Parse("2006-01-02", "2010-01-05")
	assert.NoError(t, err)
	future, err := time.Parse("2006-01-02", "2112-01-01")
	assert.NoError(t, err)

	dec31Next, err := store.NextDate(dec31)
	assert.NoError(t, err)
	jan3Next, err := store.NextDate(jan3)
	assert.NoError(t, err)
	jan4Next, err := store.NextDate(jan4)
	assert.NoError(t, err)

	assert.Equal(t, jan4, *dec31Next)
	assert.Equal(t, jan4, *jan3Next)
	assert.Equal(t, jan5, *jan4Next)

	_, err = store.NextDate(future)
	assert.Error(t, err)
}

func TestPostgresStoreGet(t *testing.T) {
	store := getDB(t)
	defer store.Close()

	jan4, err := time.Parse("2006-01-02", "2010-01-04")
	assert.NoError(t, err)
	data, err := store.Get(jan4)
	assert.NoError(t, err)

	//  ticker | open  | close | adj_close |  low  | high |   volume    |    date
	// --------+-------+-------+-----------+-------+------+-------------+------------
	//  MSFT   | 30.62 | 30.95 | 24.827723 | 30.59 | 31.1 | 3.84091e+07 | 2010-01-04

	prices := data.TickersToPrices["MSFT"]
	assert.Equal(t, 30.62, prices.Open)
	assert.Equal(t, 30.95, prices.Close)
	assert.Equal(t, 24.827723, prices.AdjClose)
	assert.Equal(t, 30.59, prices.Low)
	assert.Equal(t, 31.1, prices.High)
	assert.Equal(t, 3.84091e+07, prices.Volume)
}
