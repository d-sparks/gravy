package main

import (
	"flag"
	"log"

	gspcpipeline "github.com/d-sparks/gravy/data/dailyprices/GSPC/pipeline"
)

var (
	prices = flag.String("prices", "./data/dailyprices/GSPC/raw/^GSPC.csv", "S&P500 prices CSV.")
	dbURL  = flag.String("db", "postgres://localhost/gravy?sslmode=disable", "Postgres DB connection string.")
)

func main() {
	flag.Parse()

	if err := gspcpipeline.Pipeline(*prices, *dbURL); err != nil {
		log.Fatalf(err.Error())
	}
}
