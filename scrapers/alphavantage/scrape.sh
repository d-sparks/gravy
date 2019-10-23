#!/bin/bash

go run scrapers/alphavantage/main.go \
  --hostname="https://www.alphavantage.co/query" \
  --apikey="${ALPHAVANTAGE_API_KEY}" \
  --symbols="data/sp500_top" \
  --outputdir="data/alphavantage"
