package main

import (
	"flag"
	"log"

	dailypricespipeline "github.com/d-sparks/gravy/data/dailyprices/pipeline"
)

var (
	prices  = flag.String("prices", "./data/dailyprices/raw/historical_stock_prices.csv", "Stock prices CSV.")
	tickers = flag.String("tickers", "./data/dailyprices/raw/historical_stocks.csv", "Tickers CSV.")
	dbURL   = flag.String("db", "postgres://localhost/gravy?sslmode=disable", "Postgres DB connection string.")
)

func main() {
	flag.Parse()

	if err := dailypricespipeline.Pipeline(*prices, *tickers, *dbURL); err != nil {
		log.Fatalf(err.Error())
	}
}
