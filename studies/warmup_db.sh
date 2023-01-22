#!/bin/bash

trap "kill 0" EXIT

OUTPUT_DIR=/tmp/readonlyfizz
mkdir -p "${OUTPUT_DIR}"

# Supervisor
go run cmd/supervisor/main.go \
  > "${OUTPUT_DIR}/supervisorstdout" &

# Algorithms
go run cmd/algorithm/readonly/main.go \
  --id="readonly" \
  --port=17506 \
  > "${OUTPUT_DIR}/readonlyout" &

# Run
go run cmd/begin_backtest/main.go \
  --start="2005-02-25" \
  --end="2100-01-01" \
  --output_dir="${OUTPUT_DIR}" \
  --algorithms="readonly@localhost:17506"
