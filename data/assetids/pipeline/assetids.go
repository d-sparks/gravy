package assetidspipeline

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

// UpdateAssetIDs will load all (ticker, exchange) pairs from dailyprices and find the next available id from assetids,
// then for each pair that isn't in assetids, assigns the next id.
// TODO: This can probably be a SQL join.
func UpdateAssetIDs(postgresURL string, assetIDsTable string, dailyPricesTable string) error {
	// Connect to postgres.
	log.Printf("Connecting to database `%s`", postgresURL)
	db, err := sql.Open("postgres", postgresURL)
	if err != nil {
		return fmt.Errorf("Error connecting to postgres: %s", err.Error())
	}
	defer db.Close()

	// Get all (ticker, exchange) pairs.
	tickers := []string{}
	exchanges := []string{}
	rows, err := db.Query(fmt.Sprintf("SELECT distinct ticker, exchange FROM %s;", dailyPricesTable))
	if err != nil {
		return fmt.Errorf("Error querying for (ticker, exchange) pairs from dailyprices: %s", err.Error())
	}
	for rows.Next() {
		var ticker string
		var exchange string
		err := rows.Scan(&ticker, &exchange)
		if err != nil {
			return fmt.Errorf("Error scanning dailyprices: %s", err.Error())
		}
		tickers = append(tickers, ticker)
		exchanges = append(exchanges, exchange)
	}
	if rows.Err() != nil {
		return fmt.Errorf("Couldn't scan all ticker exchange pairs: %s", rows.Err())
	}
	rows.Close()

	// Find next available ID.
	nextID := int64(0)
	row := db.QueryRow(
		fmt.Sprintf("SELECT (CASE WHEN max(id) IS NULL THEN 0 else max(id) END) as m FROM %s;", assetIDsTable),
	)
	if err = row.Scan(&nextID); err != nil {
		return fmt.Errorf("Error scanning rows for max id: %s", err.Error())
	}
	nextID += 1

	// Record all next IDs.
	for i := 0; i < len(tickers); i++ {
		ticker := tickers[i]
		exchange := exchanges[i]

		// Check if this asset pair already has an ID.
		row = db.QueryRow(
			fmt.Sprintf(
				"SELECT count(*) FROM %s WHERE exchange=$1 and ticker=$2;",
				assetIDsTable,
			),
			exchange,
			ticker,
		)
		var count int
		if err = row.Scan(&count); err != nil {
			return fmt.Errorf("Error scanning asset ids for pair: %s", err.Error())
		}
		if count > 0 {
			continue
		}

		// Assign a new asset id to this pair.
		_, err = db.Exec(
			fmt.Sprintf(
				"INSERT INTO %s (exchange, ticker, id) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING;",
				assetIDsTable,
			),
			exchange,
			ticker,
			nextID,
		)
		if err != nil {
			return fmt.Errorf("Error inserting new ID: %s", err.Error())
		}

		// Increment nextID
		nextID++
	}

	return nil
}
