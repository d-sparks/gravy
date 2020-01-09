package main

import (
	"log"

	"github.com/d-sparks/gravy/kaggle"
	"github.com/spf13/cobra"
)

var kaggleCmd = &cobra.Command{
	Use:   "kaggle",
	Short: "trading pipeline for kaggle data",
	Run:   kaggleFn,
}

var (
	prices  string
	tickers string
	dbURL   string
)

func init() {
	pipelinesCmd.AddCommand(kaggleCmd)

	f := kaggleCmd.Flags()
	f.StringVarP(&prices, "prices", "p", "./kaggle/data/historical_stock_prices.csv", "Stock prices CSV (Kaggle).")
	f.StringVarP(&tickers, "tickers", "t", "./kaggle/data/historical_stocks.csv", "Stock symbols CSV (Kaggle)")
	f.StringVarP(&dbURL, "db", "d", "postgres://localhost/gravy?sslmode=disable", "Postgres DB connection string")
}

func kaggleFn(cmd *cobra.Command, args []string) {
	if err := kaggle.Pipeline(prices, tickers, dbURL); err != nil {
		log.Fatalf(err.Error())
	}
}
