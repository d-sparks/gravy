package dailypricespipeline

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
	"time"

	// PostGRES
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
func pricesAndDatesPipeline(filename string, db *sql.DB, pricesTable, datesTable string, exchange *Exchange) error {
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
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("Error beginning transaction in prices db: `%s`", err.Error())
	}
	dates := map[time.Time]struct{}{}
	inserted := 0
	for row, err = reader.Read(); err == nil; row, err = reader.Read() {
		if err != nil {
			return fmt.Errorf("Error parsing row in %s: %s", filename, err.Error())
		} else if len(row) != 10 {
			return fmt.Errorf("Bad row in %s: %s", filename, row)
		}
		date, err := time.Parse("20060102", row[2])
		if err != nil {
			return fmt.Errorf("Error parsing date '%s' in %s: %s", row[2], filename, err.Error())
		}
		dates[date] = struct{}{}
		query := fmt.Sprintf(
			"INSERT INTO %s (ticker, exchange, open, close, low, high, volume, date)\n",
			pricesTable,
		)
		query += "VALUES ($1, $2, $3, $4, $5, $6, $7, $8)"
		if _, err = db.Exec(
			query,
			strings.ReplaceAll(row[0], ".US", ""),
			exchange.String(),
			row[4],
			row[7],
			row[6],
			row[5],
			row[8],
			date,
		); err != nil {
			return fmt.Errorf("Error inserting row %s in %s: %s", row, filename, err.Error())
		}
		inserted++
	}
	if err != nil && err != io.EOF {
		return fmt.Errorf("Error reading prices file %s: `%s`", filename, err.Error())
	}
	log.Printf("Inserted %d prices\n", inserted)

	// Put dates into database.
	for date := range dates {
		if _, err := db.Exec(
			fmt.Sprintf(
				"INSERT INTO %s (date, %s) VALUES ($1, $2) ON CONFLICT DO NOTHING",
				datesTable,
				strings.ToLower(exchange.String()),
			),
			date,
			true,
		); err != nil {
			return fmt.Errorf("Error inserting date %s: %s", date.Format("2006-01-02"), err.Error())
		}
	}

	return tx.Commit()
}

// Pipeline populates the tickers and prices databases from files (from stooq data).
func Pipeline(pricesFolder string, dbURL string, exchange *Exchange, startAt string) error {
	// Connect to database.
	log.Printf("Connecting to database `%s`", dbURL)
	db, err := sql.Open("postgres", dbURL)
	defer db.Close()
	if err != nil {
		return fmt.Errorf("Error connecting to database: %s", err.Error())
	}

	// Get all files in the directory.
	files, err := ioutil.ReadDir(pricesFolder)
	if err != nil {
		return fmt.Errorf("Error opening prices folder: %s", err.Error())
	}

	// Process prices.
	for _, fileInfo := range files {
		if fileInfo.Name() < startAt {
			continue
		}
		filename := path.Join(pricesFolder, fileInfo.Name())
		log.Printf("Processing prices from file `%s` to table `%s`...\n", filename, PricesTable)
		if err := pricesAndDatesPipeline(filename, db, PricesTable, DatesTable, exchange); err != nil {
			return fmt.Errorf("Error processing prices: `%s`", err.Error())
		}
	}
	log.Println("Done processing prices...")

	return nil
}
