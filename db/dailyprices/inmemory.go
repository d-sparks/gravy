package dailyprices

import (
	"fmt"
	"time"

	"github.com/d-sparks/gravy/db"
)

type InMemoryStore struct {
	data map[time.Time]*db.Data
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{data: map[time.Time]*db.Data{}}
}

func NewInMemoryStoreFromFile(filename string) *InMemoryStore {
	store := InMemoryStore{data: map[time.Time]*db.Data{}}
	// TODO(dansparks): This will essentially be kaggle/pipeline.go.
	return &store
}

func (s *InMemoryStore) Get(date time.Time) (*db.Data, error) {
	data, ok := s.data[date]
	if !ok {
		return nil, fmt.Errorf("No in memory hit for `%s`", date.Format("2006-01-02"))
	}
	return data, nil
}

func (s *InMemoryStore) Set(date time.Time, data *db.Data) {
	s.data[date] = data
}
