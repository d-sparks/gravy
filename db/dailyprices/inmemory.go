package dailyprices

import (
	"fmt"
	"time"

	"github.com/d-sparks/gravy/db"
)

// Serves prices from local memory.
type InMemoryStore struct {
	data  map[time.Time]*db.Data
	dates []time.Time
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{data: map[time.Time]*db.Data{}}
}

// Not yet implemented.
func NewInMemoryStoreFromFile(filename string) *InMemoryStore {
	store := InMemoryStore{data: map[time.Time]*db.Data{}}
	// TODO(dansparks): This will essentially be kaggle/pipeline.go.
	return &store
}

// Return prices for the given date.
func (i *InMemoryStore) Get(date time.Time) (*db.Data, error) {
	data, ok := i.data[date]
	if !ok {
		return nil, fmt.Errorf("No in memory hit for `%s`", date.Format("2006-01-02"))
	}
	return data, nil
}

// Add the date to the store. Dates should be added in sequential order.
func (i *InMemoryStore) Set(date time.Time, data *db.Data) {
	i.data[date] = data
	i.dates = append(i.dates, date)
}
