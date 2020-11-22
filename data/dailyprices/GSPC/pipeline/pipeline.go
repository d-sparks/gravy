package gspcpipeline

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	dailypricespipeline "github.com/d-sparks/gravy/data/dailyprices/pipeline"
)

// Pipeline populates the tickers and prices databases from files (from kaggle data).
func Pipeline(gspcFilename string, dbURL string) error {
	// Connect to database.
	log.Printf("Connecting to database `%s`", dbURL)
	db, err := sql.Open("postgres", dbURL)
	defer db.Close()
	if err != nil {
		return fmt.Errorf("Error connecting to database: %s", err.Error())
	}

	// Process prices.
	log.Printf(
		"Processing gspc from file `%s` to table `%s`...\n",
		gspcFilename,
		dailypricespipeline.PricesTable,
	)
	if err := pricesPipeline(gspcFilename, db, dailypricespipeline.PricesTable); err != nil {
		return fmt.Errorf("Error processing gspc data: `%s`", err.Error())
	}
	log.Println("Done processing gspc...")

	return nil
}

// Populates the dailyprices table. This loads all the kaggle data into memory to group by date and then interpolate
// missing tickers prices. Also populates the dates table.
func pricesPipeline(filename string, db *sql.DB, pricesTable string) error {
	// Open file as a csv reader.
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("Error opening prices file: `%s`", err.Error())
	}
	defer file.Close()
	reader := csv.NewReader(file)

	// Skip headers.
	row, err := reader.Read()

	// Read the rows and put them in the dailyprices table.
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("Error beginning transaction in prices db: `%s`", err.Error())
	}
	dirty := 0
	inserts := 0
	errs := 0
	for row, err = reader.Read(); err == nil; row, err = reader.Read() {
		processed := dirty + inserts + errs
		if processed%5000 == 4999 {
			log.Printf(
				"... %d rows inserted [succes:%d, dirty:%d, err:%d]\n",
				processed, inserts, dirty, errs,
			)
		}
		date, err := time.Parse("2006-01-02", row[0])
		if len(row) != 7 || err != nil {
			// Dirty data.
			dirty++
			continue
		}
		query := fmt.Sprintf(
			"INSERT INTO %s (ticker, open, close, adj_close, low, high, volume, date)\n",
			pricesTable,
		)
		query += "VALUES ($1, $2, $3, $4, $5, $6, $7, $8)"
		if _, err = db.Exec(query, "^GSPC", row[1], row[4], row[5], row[3], row[2], row[6], date); err != nil {
			errs++
			continue
		}
		inserts++

	}
	if err != nil && err != io.EOF {
		return fmt.Errorf("Error reading ticks file: `%s`", err.Error())
	}

	log.Println("Total errors: ", errs)

	return tx.Commit()
}
