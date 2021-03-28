#!/bin/bash

trap "kill 0" EXIT

OUTPUT_DIR=/tmp/fizzybuzzy
mkdir -p "${OUTPUT_DIR}"

# Supervisor
go run cmd/supervisor/main.go \
  > "${OUTPUT_DIR}/supervisorstdout" &

# Algorithms
python3 algorithm/correlated_pairs/correlated_pairs.py \
  --id="correlatedpairs" \
  --port=17507 \
  --model_dir="algorithm/correlated_pairs/train/model" \
  > "${OUTPUT_DIR}/correlatedpairsstdout" &
go run cmd/algorithm/buyandhold/main.go \
  --id="buyandhold" \
  --port=17502 \
  > "${OUTPUT_DIR}/buyandholdstdout" &
go run cmd/algorithm/buyspy/main.go \
  --id="buyspy" \
  --port=17503 \
  > "${OUTPUT_DIR}/buyspystdout" &

# Run
go run main.go \
  --start="2005-02-25" \
  --end="2020-11-13" \
  --output_dir="${OUTPUT_DIR}" \
  --algorithms="correlatedpairs@localhost:17507,buyandhold@localhost:17502,buyspy@localhost:17503"
