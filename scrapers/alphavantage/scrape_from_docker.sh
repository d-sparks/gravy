#!/bin/bash

export OUTPUT_DIR="/go/src/github.com/d-sparks/ace-of-trades/scrapers/alphavantage/output"
mkdir -p $OUTPUT_DIR
/go/bin/alphavantage \
  --hostname="https://www.alphavantage.co/query" \
  --apikey="${ALPHAVANTAGE_API_KEY}" \
  --symbols="sp500_top" \
  --outputdir="${OUTPUT_DIR}"
