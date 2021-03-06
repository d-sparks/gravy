#!/bin/bash

trap "kill 0" EXIT

OUTPUT_DIR=/tmp/fizz
mkdir -p "${OUTPUT_DIR}"

# Supervisor
go run cmd/supervisor/main.go \
  > "${OUTPUT_DIR}/supervisorstdout" &

# Algorithms
go run cmd/algorithm/buyandhold/main.go \
  --id="buyandhold" \
  --port=17502 \
  > "${OUTPUT_DIR}/buyandholdstdout" &
go run cmd/algorithm/buyspy/main.go \
  --id="buyspy" \
  --port=17503 \
  > "${OUTPUT_DIR}/buyspystdout" &

# Run
go run cmd/begin_backtest/main.go \
  --start="2005-02-25" \
  --end="2007-02-25" \
  --output_dir="${OUTPUT_DIR}" \
  --algorithms="buyandhold@localhost:17502,buyspy@localhost:17503"
