package main

import (
	"flag"
	"log"

	dailypricespipeline "github.com/d-sparks/gravy/data/dailyprices/pipeline"
)

var (
	pricesFolder = flag.String(
		"prices_folder",
		"./data/dailyprices/raw/daily/us/nyse stocks/1",
		"Stock prices folder containing csvs.",
	)
	dbURL       = flag.String("db", "postgres://localhost/gravy?sslmode=disable", "Postgres DB connection string.")
	exchangeStr = flag.String("exchange", "NYSE", "String representation of the exchange for these csvs.")
	startAt     = flag.String("start_at", "", "Skip up until the specified file.")
)

func main() {
	flag.Parse()

	exchange, err := dailypricespipeline.ParseExchange(*exchangeStr)
	if err != nil {
		log.Fatalf(err.Error())
	}

	if err := dailypricespipeline.Pipeline(*pricesFolder, *dbURL, exchange, *startAt); err != nil {
		log.Fatalf(err.Error())
	}
}
