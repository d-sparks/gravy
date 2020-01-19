package dailyprices

import (
	"fmt"
	"sort"
	"time"

	"github.com/d-sparks/gravy/db"
)

// InMemoryStore serves prices from local memory.
type InMemoryStore struct {
	data                map[time.Time]*db.Data
	dates               []time.Time
	nextDates           map[time.Time]time.Time
	nextDatesCalculated bool
}

// NewInMemoryStore returns a new InMemoryStore with empty data.
func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		data:                map[time.Time]*db.Data{},
		dates:               []time.Time{},
		nextDates:           map[time.Time]time.Time{},
		nextDatesCalculated: false,
	}
}

// NewInMemoryStoreFromFile is not yet implemented.
func NewInMemoryStoreFromFile(filename string) *InMemoryStore {
	store := NewInMemoryStore()
	// TODO(dansparks): This will essentially be kaggle/pipeline.go.
	return store
}

// ValidDate implements db interface and returns whether the date is valid.
func (i *InMemoryStore) ValidDate(date time.Time) (bool, error) {
	_, ok := i.data[date]
	return ok, nil
}

// NextDate returns the next date. Calls calculate next dates if it hasn't been called.
func (i *InMemoryStore) NextDate(date time.Time) (*time.Time, error) {
	if valid, err := i.ValidDate(date); err != nil {
		return nil, err
	} else if valid {
		if !i.nextDatesCalculated {
			i.calculateNextDates()
		}
		nextDate, ok := i.nextDates[date]
		if !ok {
			return nil, fmt.Errorf("No next date")
		}
		return &nextDate, nil
	}
	// Brute force!
	var nextDate *time.Time = nil
	for _, otherDate := range i.dates {
		if nextDate == nil || (otherDate.After(date) && nextDate.After(otherDate)) {
			nextDate = &otherDate
		}
	}
	if nextDate == nil {
		return nil, fmt.Errorf("No next date")
	}
	return nextDate, nil
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
	i.nextDatesCalculated = false
}

// calculateNextDates calcultes the date after each date. Not part of db interface.
func (i *InMemoryStore) calculateNextDates() {
	i.nextDates = map[time.Time]time.Time{}
	dateStrs := sort.StringSlice{}
	parsed := map[string]time.Time{}
	for _, date := range i.dates {
		dateStr := date.Format("2006-01-02")
		dateStrs = append(dateStrs, dateStr)
		parsed[dateStr] = date
	}
	dateStrs.Sort()
	for ix := 1; ix < len(dateStrs); ix++ {
		i.nextDates[parsed[dateStrs[ix-1]]] = parsed[dateStrs[ix]]
	}
	i.nextDatesCalculated = true
}
