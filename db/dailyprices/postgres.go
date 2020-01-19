package dailyprices

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/Clever/go-utils/stringset"
	"github.com/d-sparks/gravy/db"
	_ "github.com/lib/pq" // postgres for sql
)

// PostgresStore serves daily prices from a postgres database.
type PostgresStore struct {
	db          *sql.DB
	pricesTable string
	datesTable  string
}

// NewPostgresStore creates a new postgres store pointing at a given db and tables.
func NewPostgresStore(dbURL string, pricesTable, datesTable string) (*PostgresStore, error) {
	// Connect to database.
	log.Printf("Connecting to database `%s`", dbURL)
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, fmt.Errorf("Error connecting to database: `%s`", err.Error())
	}

	return &PostgresStore{db, pricesTable, datesTable}, nil
}

// Close the database connection.
func (s *PostgresStore) Close() {
	s.db.Close()
}

// ValidDate returns whether the date is a trading date, and is part of the db interface.
func (s *PostgresStore) ValidDate(date time.Time) (bool, error) {
	rows, err := s.db.Query("SELECT date WHERE date = $1", date.Format("2006-01-02"))
	if err != nil {
		return false, fmt.Errorf("Error reading from db: `%s`", err.Error())
	}
	return rows.Next(), nil
}

// NextDate returns the next trading date after the given date, assuming given date is valid.
func (s *PostgresStore) NextDate(date time.Time) (*time.Time, error) {
	rows, err := s.db.Query("SELECT date WHERE date > $1 ORDER BY date DESC limit 1", date.Format("2006-01-02"))
	if err != nil {
		return nil, fmt.Errorf("Error reading from db: `%s`", err.Error())
	}
	if !rows.Next() {
		return nil, fmt.Errorf("No dates found")
	}
	var nextDate time.Time
	if err = rows.Scan(&nextDate); err != nil {
		return nil, fmt.Errorf("Error parsing nextDate: `%s`", err.Error())
	}
	return &nextDate, nil
}

// Get proces for a specific date.
func (s *PostgresStore) Get(date time.Time) (*db.Data, error) {
	// Query database.
	rows, err := s.db.Query(
		fmt.Sprintf(
			"SELECT ticker, open, close, adj_close, low, high, volume FROM %s WHERE date = $1",
			s.pricesTable,
		),
		date.Format("2006-01-02"),
	)
	if err != nil {
		return nil, fmt.Errorf("Error reading from db: `%s`", err.Error())
	}

	// Construct window.
	data := db.Data{TickersToPrices: map[string]db.Prices{}, Tickers: stringset.New()}
	for rows.Next() {
		var ticker string
		// TODO(dansparks): Use a db.Prices here.
		var open, cloze, adjClose, low, high, volume float64
		if err = rows.Scan(&ticker, &open, &cloze, &adjClose, &low, &high, &volume); err != nil {
			return nil, fmt.Errorf("Error while parsing row: `%s`", err.Error())
		}
		data.TickersToPrices[ticker] = db.Prices{
			Open:     open,
			Close:    cloze,
			AdjClose: adjClose,
			Low:      low,
			High:     high,
			Volume:   volume,
		}
		data.Tickers.Add(ticker)
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("Error constructing response: `%s`", rows.Err().Error())
	}

	return &data, nil
}

// NextDate is for the db interface, returns the next valid date.

// AllDates returns distinct dates in the database.
func (s *PostgresStore) AllDates() ([]time.Time, error) {
	// Query for dates.
	rows, err := s.db.Query(fmt.Sprintf("SELECT DISTINCT date FROM %s ORDER BY date", s.datesTable))
	if err != nil {
		return nil, fmt.Errorf("Error querying for distinct dates: `%s`", err.Error())
	}

	// Scan and parse dates into a slice.
	dates := []time.Time{}
	for rows.Next() {
		var dateStr string
		if err := rows.Scan(&dateStr); err != nil {
			return nil, fmt.Errorf("Error scanning date `%s` from postgres: `%s`", dateStr, err.Error())
		}
		date, err := time.Parse("2006-01-02T15:04:05Z", dateStr)
		if err != nil {
			return nil, fmt.Errorf("Could not parse date `%s`: `%s`", dateStr, err.Error())
		}
		dates = append(dates, date)
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("Error scanning rows for distinct dates: `%s`", rows.Err().Error())
	}

	return dates, nil
}
