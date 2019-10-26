package dailywindow

import (
	"encoding/json"
	"time"

	"github.com/d-sparks/gravy/db"
	"github.com/d-sparks/gravy/gravyutil"
	"github.com/d-sparks/gravy/trading"
)

const Name = "dailywindow"

type InMemoryStore struct {
	data  map[time.Time]*trading.Window
	dates []time.Time
}

func NewInMemoryStore(filename string) *InMemoryStore {
	store := InMemoryStore{data: map[time.Time]*trading.Window{}}
	scanner := gravyutil.FileScannerOrDie(filename)
	for scanner.Scan() {
		window := trading.Window{}
		gravyutil.FatalIfErr(json.Unmarshal(scanner.Bytes(), &window))
		store.data[window.Begin] = &window
		store.dates = append(store.dates, window.Begin)
	}
	gravyutil.FatalIfErr(scanner.Err())
	return &store
}

func (s *InMemoryStore) Get(date time.Time) db.Data {
	return db.Data{Window: *s.data[date]}
}

func (s *InMemoryStore) Dates() []time.Time {
	return s.dates
}
