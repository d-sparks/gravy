package dailywindow

import (
	"time"

	"github.com/d-sparks/gravy/db"
	"github.com/d-sparks/gravy/trading"
)

const Name = "dailywindow"

type InMemoryStore struct {
	data  map[time.Time]trading.Window
	dates []time.Time
}

func NewInMemoryStore(filename string) InMemoryStore {
	return InMemoryStore{data: map[time.Time]trading.Window{}}
	// TODO: parse.
}

func (s *InMemoryStore) Get(date time.Time) db.Data {
	return db.Data{Window: s.data[date]}
}

func (s *InMemoryStore) Dates() []time.Time {
	return s.dates
}
