package dailywindow

import (
	"encoding/json"
	"log"
	"time"

	"github.com/Clever/go-utils/stringset"
	"github.com/d-sparks/gravy/db"
	"github.com/d-sparks/gravy/gravyutil"
	"github.com/d-sparks/gravy/trading"
)

const Name = "dailywindow"

type InMemoryStore struct {
	data  map[time.Time]*trading.Window
	dates []time.Time
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{data: map[time.Time]*trading.Window{}, dates: []time.Time{}}
}

func NewInMemoryStoreFromFile(filename string) *InMemoryStore {
	store := InMemoryStore{data: map[time.Time]*trading.Window{}}
	scanner := gravyutil.FileScannerOrDie(filename)
	log.Println("Beginning to load in memory dailywindow store...")
	for scanner.Scan() {
		window := trading.Window{}
		gravyutil.FatalIfErr(json.Unmarshal(scanner.Bytes(), &window))
		// TODO(dansparks): Maybe serialize these in the csv.
		window.Symbols = stringset.New()
		for symbol, _ := range window.Open {
			window.Symbols.Add(symbol)
		}
		store.data[window.Begin] = &window
		store.dates = append(store.dates, window.Begin)
	}
	gravyutil.FatalIfErr(scanner.Err())
	log.Printf("Loaded %d dailywindows into memory...", len(store.dates))
	return &store
}

func (s *InMemoryStore) Get(date time.Time) db.Data {
	return db.Data{Window: *s.data[date]}
}

func (s *InMemoryStore) Set(date time.Time, window *trading.Window) {
	s.dates = append(s.dates, date)
	s.data[date] = window
}

func (s *InMemoryStore) Dates() []time.Time {
	return s.dates
}
