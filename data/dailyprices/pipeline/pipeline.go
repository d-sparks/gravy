package dailypricespipeline

import (
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"reflect"
	"strings"
	"sync"
	"time"

	// PostGRES
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const (
	// PricesTable holds the daily prices.
	PricesTable = "dailyprices"

	// DatesTable holds the valid trading dates.
	DatesTable = "tradingdates"
)

// Exchange type (for recording trading dates).
type Exchange int

const (
	// NYSE is the NYSE
	NYSE Exchange = iota
	// NYSE_ETF is the NYSEF ETF's
	NYSE_ETF
	// NASDAQ is stocks on the NASDAQ
	NASDAQ
	// NASDAQ_ETF is the NASDAQ's ETF's
	NASDAQ_ETF
	// NYSEMKT is stocks on the American Stock Exchange.
	NYSEMKT
	// NYSEMKT_ETF is ETFs on the American Stock Exchange.
	NYSEMKT_ETF
)

var exchangeStrings = map[Exchange]string{
	NYSE:        "NYSE",
	NYSE_ETF:    "NYSE_ETF",
	NASDAQ:      "NASDAQ",
	NASDAQ_ETF:  "NASDAQ_ETF",
	NYSEMKT:     "NYSEMKT",
	NYSEMKT_ETF: "NYSEMKT_ETF",
}

// ExchangeString is the string for which each exchange is stored in the database.
func (e *Exchange) String() string {
	return exchangeStrings[*e]
}

// ParseExchange parses a string into an exchange type or returns an error.
func ParseExchange(exchange string) (*Exchange, error) {
	for k, v := range exchangeStrings {
		if v == exchange {
			return &k, nil
		}
	}
	return nil, fmt.Errorf("Unknown exchange: %s", exchange)
}

// Populates the dailyprices table. This loads all the kaggle data into memory to group by date and then interpolate
// missing tickers prices. Also populates the dates table.
func pricesAndDatesPipeline(filename string, db *sqlx.DB, pricesTable, datesTable string, exchange *Exchange) error {
	// Open file as a csv reader.
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("Error opening prices file: `%s`", err.Error())
	}
	defer file.Close()
	reader := csv.NewReader(file)

	// Verify headers.
	row, err := reader.Read()
	if err == io.EOF {
		return nil
	} else if err != nil {
		return fmt.Errorf("Error reading csv headers in file %s: %s", filename, err.Error())
	}
	if row[0] != "<TICKER>" || row[1] != "<PER>" || row[2] != "<DATE>" || row[3] != "<TIME>" ||
		row[4] != "<OPEN>" || row[5] != "<HIGH>" || row[6] != "<LOW>" || row[7] != "<CLOSE>" ||
		row[8] != "<VOL>" || row[9] != "<OPENINT>" {

		return fmt.Errorf("Unexpected headers in %s", filename)
	}

	// Read the rows and put them in the dailyprices table.
	tx, err := db.Beginx()
	if err != nil {
		return fmt.Errorf("Error beginning transaction in prices db: `%s`", err.Error())
	}
	dates := map[time.Time]struct{}{}
	inserted := 0

	query := fmt.Sprintf(
		`INSERT INTO %s (ticker, exchange, open, close, low, high, volume, date)
			 VALUES (:ticker, :exchange, :open, :close, :low, :high, :volume, :date) ON CONFLICT DO NOTHING`,
		pricesTable,
	)

	type Entry struct {
		Ticker   string    `db:"ticker"`
		Exchange string    `db:"exchange"`
		Open     string    `db:"open"`
		Close    string    `db:"close"`
		Low      string    `db:"low"`
		High     string    `db:"high"`
		Volume   string    `db:"volume"`
		Date     time.Time `db:"date"`
	}

	// TODO(desa): is changing this to `ON CONFLICT DO NOTHING` ok?
	dateQuery := fmt.Sprintf(
		"INSERT INTO %s (date, %s) VALUES (:date, :one) ON CONFLICT DO NOTHING",
		datesTable,
		strings.ToLower(exchange.String()),
	)
	type DateEntry struct {
		Date time.Time `db:"date"`
		One  bool      `db:"one"`
		Two  bool
	}

	entryV := reflect.ValueOf(Entry{})
	dateEntryV := reflect.ValueOf(DateEntry{})
	maxEntryBatch := ((1 << 16) / entryV.NumField()) - 1
	maxDateBatch := ((1 << 16) / dateEntryV.NumField()) - 1

	entries := make([]Entry, 0, maxEntryBatch)
	dateEntries := make([]DateEntry, 0, maxDateBatch)
	for row, err = reader.Read(); err == nil; row, err = reader.Read() {
		if err != nil {
			return fmt.Errorf("Error parsing row in %s: %s", filename, err.Error())
		}

		if len(row) != 10 {
			return fmt.Errorf("Bad row in %s: %s", filename, row)
		}

		date, err := time.Parse("20060102", row[2])
		if err != nil {
			return fmt.Errorf("Error parsing date '%s' in %s: %s", row[2], filename, err.Error())
		}

		if _, ok := dates[date]; !ok {
			dates[date] = struct{}{}
			dateEntries = append(dateEntries, DateEntry{
				Date: date,
				One:  true,
				Two:  true,
			})
		}

		entries = append(entries, Entry{
			Ticker:   strings.ReplaceAll(row[0], ".US", ""),
			Exchange: exchange.String(),
			Open:     row[4],
			Close:    row[7],
			Low:      row[6],
			High:     row[5],
			Volume:   row[8],
			Date:     date,
		})
		inserted++

		if err != nil && err != io.EOF {
			return fmt.Errorf("Error reading prices file %s: `%s`", filename, err.Error())
		}

		if len(entries) == maxEntryBatch {
			if _, err = tx.NamedExec(query, entries); err != nil {
				return fmt.Errorf("Error inserting row %s in %s: %s", row, filename, err.Error())
			}
			log.Printf("Inserted %d prices\n", len(entries))
			entries = entries[:0]
		}

		if len(dateEntries) == maxDateBatch {
			if _, err := tx.NamedExec(dateQuery, dateEntries); err != nil {
				return fmt.Errorf("Error inserting dates: %s", err.Error())
			}
			log.Printf("Inserted %d dates\n", len(dateEntries))
			dateEntries = dateEntries[:0]
		}
	}

	if len(entries) > 0 {
		if _, err = tx.NamedExec(query, entries); err != nil {
			return fmt.Errorf("Error inserting row %s in %s: %s", row, filename, err.Error())
		}
		log.Printf("Inserted %d prices\n", len(entries))
	}

	if len(dateEntries) > 0 {
		if _, err := tx.NamedExec(dateQuery, dateEntries); err != nil {
			return fmt.Errorf("Error inserting dates: %s", err.Error())
		}
		log.Printf("Inserted %d dates\n", len(dateEntries))
	}

	return tx.Commit()
}

type Pipeline struct {
	db        *sqlx.DB
	exchange  *Exchange
	startAt   time.Time
	pricesDir string
}

func (p *Pipeline) Exec() error {
	// Get all files in the directory.
	files, err := ioutil.ReadDir(p.pricesDir)
	if err != nil {
		return fmt.Errorf("Error opening prices folder: %s", err.Error())
	}
	var wg sync.WaitGroup
	limiter := make(chan bool, 20)

	// Process prices.
	for _, fileInfo := range files {
		wg.Add(1)
		limiter <- true
		filename := path.Join(p.pricesDir, fileInfo.Name())
		go func(filename string) {
			defer func() {
				<-limiter
				wg.Done()
			}()
			log.Printf("Processing prices from file `%s` to table `%s`...\n", filename, PricesTable)
			if err := pricesAndDatesPipeline(filename, p.db, PricesTable, DatesTable, p.exchange); err != nil {
				log.Fatalf("Error processing prices: `%s`", err.Error())
			}
		}(filename)
	}
	wg.Wait()
	log.Println("Done processing prices...")

	return p.Close()
}

func (p *Pipeline) Close() error {
	if err := p.db.Close(); err != nil {
		return err
	}

	return nil
}

// NewPipeline populates the tickers and prices databases from files (from stooq data).
func NewPipeline(pricesFolder string, dbURL string, exchange *Exchange, startAt string) (*Pipeline, error) {
	// Connect to database.
	log.Printf("Connecting to database `%s`", dbURL)
	db, err := sqlx.Open("postgres", dbURL)
	if err != nil {
		return nil, fmt.Errorf("Error connecting to database: %s", err.Error())
	}

	return &Pipeline{
		db:        db,
		exchange:  exchange,
		startAt:   time.Now(), // TODO(desa): is this needed?
		pricesDir: pricesFolder,
	}, nil
}
