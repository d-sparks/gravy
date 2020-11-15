package dailypricespipeline

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strconv"

	"github.com/Clever/go-utils/stringset"
	// PostGRES
	_ "github.com/lib/pq"
)

const (
	pricesTable  = "dailyprices"
	tickersTable = "tickers"
	datesTable   = "tradingdates"
)

// Pipeline populates the tickers and prices databases from files (from kaggle data).
func Pipeline(pricesFilename, tickersFilename, dbURL string) error {
	// Connect to database.
	log.Printf("Connecting to database `%s`", dbURL)
	db, err := sql.Open("postgres", dbURL)
	defer db.Close()
	if err != nil {
		return fmt.Errorf("Error connecting to database: %s", err.Error())
	}

	// Process tickers.
	log.Printf("Processing tickers from file `%s` to table `%s`...\n", tickersFilename, tickersTable)
	if err := tickersPipeline(tickersFilename, db, tickersTable); err != nil {
		return fmt.Errorf("Error processing tickers: `%s`", err.Error())
	}
	log.Println("Done processing tickers...")

	// Process prices.
	log.Printf("Processing prices from file `%s` to table `%s`...\n", pricesFilename, pricesTable)
	if err := pricesAndDatesPipeline(pricesFilename, db, pricesTable, datesTable); err != nil {
		return fmt.Errorf("Error processing ticks: `%s`", err.Error())
	}
	log.Println("Done processing prices...")

	return nil
}

// Populates the given database with ticker data.
func tickersPipeline(filename string, db *sql.DB, table string) error {
	// Open file as a csv reader.
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("Error opening tickers file: `%s`", err.Error())
	}
	defer file.Close()
	reader := csv.NewReader(file)

	// Skip headers.
	row, err := reader.Read()

	// Process rows into database.
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("Error beginning transaction in tickers db: `%s`", err.Error())
	}
	for row, err = reader.Read(); err == nil; row, err = reader.Read() {
		if len(row) != 5 {
			return fmt.Errorf("Expected tickers row with 5 columns: `%s`", row)
		}
		if len(row[0]) > 5 {
			// Dirty data.
			continue
		}
		query := fmt.Sprintf("INSERT INTO %s (ticker, exchange, name, sector, industry)\n", table)
		query += "VALUES ($1, $2, $3, $4, $5)"
		_, err := db.Exec(query, row[0], row[1], row[2], row[3], row[4])
		if err != nil {
			return fmt.Errorf("Error inserting ticker row `%s`: `%s`", row, err.Error())
		}

	}
	if err != nil && err != io.EOF {
		return fmt.Errorf("Error processing ticker rows: `%s`", err.Error())
	}
	return tx.Commit()
}

// Populates the dailyprices table. This loads all the kaggle data into memory to group by date and then interpolate
// missing tickers prices. Also populates the dates table.
func pricesAndDatesPipeline(filename string, db *sql.DB, pricesTable, datesTable string) error {
	// Open file as a csv reader.
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("Error opening prices file: `%s`", err.Error())
	}
	defer file.Close()
	reader := csv.NewReader(file)

	// Skip headers.
	row, err := reader.Read()

	// Group by date and ticker, storing open, close, adj_close, low, high, and volume as float64.
	dates := sort.StringSlice{}
	datesSet := stringset.New()
	dateToTickerToCols := map[string]map[string][]float64{}
	dateToTickers := map[string]stringset.StringSet{}
	for row, err = reader.Read(); err == nil; row, err = reader.Read() {
		if len(row) != 8 {
			return fmt.Errorf("Invalid row `%s` in prices csv: `%s`", row, err.Error())
		}

		date := row[7]
		ticker := row[0]

		// Initialize maps if necessary.
		if _, ok := dateToTickerToCols[date]; !ok {
			dateToTickerToCols[date] = map[string][]float64{}
			dateToTickers[date] = stringset.New()
			dates = append(dates, date)
		}
		dateToTickerToCols[date][ticker] = make([]float64, 6)

		// Populate columns from row.
		for i := 1; i <= 6; i++ {
			dateToTickerToCols[date][ticker][i-1], err = strconv.ParseFloat(row[i], 64)
			if err != nil {
				return fmt.Errorf("Error parsing row `%s`: `%s`", row, err.Error())
			}
		}

		// Insert into dateset (distinct dates).
		datesSet.Add(date)
	}
	if err != nil && err != io.EOF {
		return fmt.Errorf("Error reading ticks file: `%s`", err.Error())
	}

	// Interpolate data for tickers that were missing in some dates.
	dates.Sort()
	InterpolateData(dates, dateToTickerToCols, dateToTickers)

	// Put prices into database.
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("Error beginning transaction in prices db: `%s`", err.Error())
	}
	total := 0
	dirty := 0
	inserts := 0
	errs := 0
	for _, tickerToCols := range dateToTickerToCols {
		total += len(tickerToCols)
	}
	for date, tickerToCols := range dateToTickerToCols {
		for ticker, cols := range tickerToCols {
			processed := dirty + inserts + errs
			if processed%50000 == 49999 {
				progress := 100.0 * float64(processed) / float64(total)
				log.Printf(
					"... %f%% complete (%d/%d rows inserted) [succes:%d, dirty:%d, err:%d]\n",
					progress, processed, total, inserts, dirty, errs,
				)
			}
			if len(cols) != 6 || len(ticker) > 5 {
				// Dirty data.
				dirty++
				continue
			}
			query := fmt.Sprintf(
				"INSERT INTO %s (ticker, open, close, adj_close, low, high, volume, date)\n",
				pricesTable,
			)
			query += "VALUES ($1, $2, $3, $4, $5, $6, $7, $8)"
			_, err := db.Exec(query, ticker, cols[0], cols[1], cols[2], cols[3], cols[4], cols[5], date)
			if err != nil {
				errs++
				continue
			}
			inserts++
		}

	}
	log.Println("Total errors: ", errs)

	// Put dates into database.
	for date, _ := range datesSet {
		_, err := db.Exec(fmt.Sprintf("INSERT INTO %s (date) VALUES ($1)", datesTable), date)
		if err != nil {
			errs++
			continue
		}
	}

	return tx.Commit()
}

// For days when we don't have data on a specific stock, but we do have data in the previous and
// next windows, we interpolate using the mean.
func InterpolateData(
	dates []string,
	dateToTickerToCols map[string]map[string][]float64,
	dateToTickers map[string]stringset.StringSet,
) {
	// Get the last index of when a ticker was observed.
	lastListingIx := map[string]int{}
	for i := 0; i < len(dates); i++ {
		for ticker, _ := range dateToTickers[dates[i]] {
			lastListingIx[ticker] = i
		}
	}

	for i := 1; i+1 < len(dates); i++ {
		tickersPrev := dateToTickers[dates[i-1]]
		tickers := dateToTickers[dates[i]]
		missing := tickersPrev.Minus(tickers)

		for ticker, _ := range missing {
			// No need to do anything for unlisted tickers.
			if i >= lastListingIx[ticker] {
				continue
			}

			// Find the next index for which this ticker is listed.
			j := i + 1
			for ; !dateToTickers[dates[j]].Contains(ticker); j++ {
			}

			// So we have listings at index i-1 and j, and we need to fill [i, j-1].
			InterpolateForSymbol(ticker, i-1, j, dates, dateToTickerToCols)
			dateToTickers[dates[i]].Add(ticker)
		}
	}
}

// Assumes the stock is published at dates[lb] and dates[ub] but not between.
func InterpolateForSymbol(
	ticker string,
	lb int,
	ub int,
	dates []string,
	dateToTickerToCols map[string]map[string][]float64,
) error {
	// If the stock was unlisted for more than 90 days, assume the new ticker is a different
	// company.
	if ub-lb > 90 {
		return nil
	}

	// Back and front are slices of float64's, corresponding to this ticker's numeric columns at the date
	// corresponding to lb and the date corresponding to ub, respectively.
	back := dateToTickerToCols[dates[lb]][ticker]
	front := dateToTickerToCols[dates[ub]][ticker]
	if len(back) != len(front) {
		return fmt.Errorf(
			"Malformed columns for dates `%s`, `%s`, and ticker `%s`",
			dates[lb],
			dates[ub],
			ticker,
		)
	}

	// The stock wasn't published on the dates from lb to ub, so make empty column vectors for these dates.
	for dateIx := lb + 1; dateIx < ub; dateIx++ {
		dateToTickerToCols[dates[dateIx]][ticker] = make([]float64, len(back))
	}

	// Linearly interpolate based on y = y1 + m(x-x1) = ((y2-y1)/(x2-x1))*(x-x1). Note that here x is the date index
	// and y is the column (e.g. open, close, adj_close, low, high, volume) being interpolated.
	dx := float64(ub - lb)
	x1 := float64(lb)
	for colIx := 0; colIx < len(back); colIx++ {
		m := (front[colIx] - back[colIx]) / dx
		y1 := back[colIx]
		for dateIx := lb + 1; dateIx < ub; dateIx++ {
			x := float64(dateIx)
			dateToTickerToCols[dates[dateIx]][ticker][colIx] = y1 + m*(x-x1)
		}
	}

	return nil
}
