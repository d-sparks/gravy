#!/bin/bash

trap "kill 0" EXIT

OUTPUT_DIR=/tmp/foo
mkdir -p "${OUTPUT_DIR}"

# Supervisor
go run cmd/supervisor/main.go \
  > "${OUTPUT_DIR}/supervisorstdout" &

# Algorithms
python3 \
  algorithm/headsortails/heads_or_tails.py \
  --id="headsortails" \
  --port="17506" \
  --model_dir="algorithm/headsortails/train/model" \
  > "${OUTPUT_DIR}/headsortailsstdout" &
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
  --end="2020-02-02" \
  --output_dir="${OUTPUT_DIR}" \
  --algorithms="headsortails@localhost:17506,buyandhold@localhost:17502,buyspy@localhost:17503"
