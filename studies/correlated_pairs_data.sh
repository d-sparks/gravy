#!/bin/bash

trap "kill 0" EXIT

OUTPUT_DIR=/tmp/fizzbuzz
mkdir -p "${OUTPUT_DIR}"

# Supervisor
go run cmd/supervisor/main.go \
  > "${OUTPUT_DIR}/supervisorstdout" &

# Algorithms
python3 algorithm/correlated_pairs/correlated_pairs.py \
  --id="correlatedpairs" \
  --port=17507 \
  --export_training_data=true \
  --training_data_path="${OUTPUT_DIR}/correlatedpairs_data.csv" &
  > "${OUTPUT_DIR}/correlatedpairsstdout"

# Run
go run cmd/begin_backtest/main.go \
  --start="2005-02-25" \
  --end="2006-02-25" \
  --output_dir="${OUTPUT_DIR}" \
  --algorithms="correlatedpairs@localhost:17507"
