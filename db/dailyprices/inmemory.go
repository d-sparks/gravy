package dailyprices

import (
	"fmt"
	"time"

	"github.com/d-sparks/gravy/db"
)

// InMemoryStore serves prices from local memory.
type InMemoryStore struct {
	data  map[time.Time]*db.Data
	dates []time.Time
}

// NewInMemoryStore returns a new InMemoryStore with empty data.
func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{data: map[time.Time]*db.Data{}}
}

// NewInMemoryStoreFromFile is not yet implemented.
func NewInMemoryStoreFromFile(filename string) *InMemoryStore {
	store := InMemoryStore{data: map[time.Time]*db.Data{}}
	// TODO(dansparks): This will essentially be kaggle/pipeline.go.
	return &store
}

// Get returns prices for the given date.
func (i *InMemoryStore) Get(date time.Time) (*db.Data, error) {
	data, ok := i.data[date]
	if !ok {
		return nil, fmt.Errorf("No in memory hit for `%s`", date.Format("2006-01-02"))
	}
	return data, nil
}

// Set adds the date to the store. Dates should be added in sequential order.
func (i *InMemoryStore) Set(date time.Time, data *db.Data) {
	i.data[date] = data
	i.dates = append(i.dates, date)
}
