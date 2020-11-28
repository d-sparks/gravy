#!/bin/bash

trap "kill 0" EXIT

OUTPUT_DIR=/tmp/fizz
mkdir -p "${OUTPUT_DIR}"

# Supervisor
go run cmd/supervisor/main.go \
  > "${OUTPUT_DIR}/supervisorstdout" &

# Algorithms
python3 \
   algorithm/headsortails/heads_or_tails.py &

# go run cmd/algorithm/headsortails/main.go \
#   --id="headsortails" \
#   --port=17505 \
#   --mode="train" \
#   --sample_ratio=0.02 \
#   --output="${OUTPUT_DIR}/headsortails_data.csv" \
#   > "${OUTPUT_DIR}/headsortailsstdout" &

# Run
go run main.go \
  --start="2005-02-25" \
  --end="2005-03-02" \
  --output_dir="/tmp/fizz" \
  --algorithms="headsortails@localhost:17506"
