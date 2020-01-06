package dailyprices

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/d-sparks/gravy/db"
	_ "github.com/lib/pq"
)

// Serve daily prices from a postgres database.
type PostgresStore struct {
	db    *sql.DB
	table string
}

func NewPostgresStore(dbURL string, table string) (*PostgresStore, error) {
	// Connect to database.
	log.Printf("Connecting to database `%s`", dbURL)
	db, err := sql.Open("postgres", dbURL)
	defer db.Close()
	if err != nil {
		return nil, fmt.Errorf("Error connecting to database.")
	}

	return &PostgresStore{db, table}, nil
}

func (s *PostgresStore) Get(date time.Time) (*db.Data, error) {
	// Query database.
	rows, err := s.db.Query(
		"SELECT ticker, open, close, adj_close, low, high, volume from $1 WHERE date = '$2'",
		s.table,
		date.Format("2006-01-02"),
	)
	if err != nil {
		return nil, fmt.Errorf("Error reading from db: `%s`", err.Error())
	}

	// Construct window.
	data := db.Data{}
	for rows.Next() {
		var ticker string
		var open, cloze, adjClose, low, high, volume float64
		if err = rows.Scan(&ticker, &open, &adjClose, &low, &high, &volume); err != nil {
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
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("Error constructing response: `%s`", rows.Err().Error())
	}

	return &data, nil
}
