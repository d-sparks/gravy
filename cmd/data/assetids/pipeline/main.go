package main

import (
	"flag"
	"log"

	assetids "github.com/d-sparks/gravy/data/assetids/pipeline"
)

var (
	dbURL = flag.String(
		"db",
		"postgres://localhost/gravy?sslmode=disable",
		"Postgres DB connection string.",
	)
	dailyPricesTable = flag.String("dailyprices_table", "dailyprices", "Daily prices table.")
	assetIDsTable    = flag.String("assetids_table", "assetids", "Asset IDs table.")
)

func main() {
	flag.Parse()
	if err := assetids.UpdateAssetIDs(*dbURL, *assetIDsTable, *dailyPricesTable); err != nil {
		log.Fatalf(err.Error())
	}
}
