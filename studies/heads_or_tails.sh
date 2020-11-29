#!/bin/bash

trap "kill 0" EXIT

OUTPUT_DIR=/tmp/fizz
mkdir -p "${OUTPUT_DIR}"

# Supervisor
go run cmd/supervisor/main.go \
  > "${OUTPUT_DIR}/supervisorstdout" &

# Algorithms
python3 \
  algorithm/headsortails/heads_or_tails.py \
  --id="headsortails" \
  --port="17506" \
  --model_dir="algorithm/headsortails/train/model" &

# Run
go run main.go \
  --start="2005-02-25" \
  --end="2005-04-02" \
  --output_dir="/tmp/fizz" \
  --algorithms="headsortails@localhost:17506"
