#!/bin/bash

# Golang
protoc --go_out=$GOPATH/src --go-grpc_out=$GOPATH/src supervisor/proto/supervisor.proto
protoc --go_out=$GOPATH/src --go-grpc_out=$GOPATH/src data/dailyprices/proto/daily_prices.proto
protoc --go_out=$GOPATH/src --go-grpc_out=$GOPATH/src algorithm/proto/algorithm_io.proto

# Python
python3 \
  -m grpc_tools.protoc \
  --proto_path=data/dailyprices/proto/ \
  --python_out=data/dailyprices/proto/ \
  --grpc_python_out=data/dailyprices/proto/ \
  daily_prices.proto
sed \
  -i '' 's/import daily_prices_pb2 as daily__prices__pb/from . import daily_prices_pb2 as daily__prices__pb/g' \
  data/dailyprices/proto/daily_prices_pb2_grpc.py

python3 \
  -m grpc_tools.protoc \
  --proto_path=supervisor/proto/ \
  --python_out=supervisor/proto/ \
  --grpc_python_out=supervisor/proto/ \
  supervisor.proto
sed \
  -i '' 's/import supervisor_pb2 as supervisor__pb2/from . import supervisor_pb2 as supervisor__pb2/g' \
  supervisor/proto/supervisor_pb2_grpc.py

python3 \
  -m grpc_tools.protoc \
  --proto_path=algorithm/proto \
  --python_out=algorithm/proto/ \
  --grpc_python_out=algorithm/proto/ \
  algorithm_io.proto
sed \
  -i '' 's/import algorithm_io_pb2 as algorithm__io__pb2/from . import algorithm_io_pb2 as algorithm__io__pb2/g' \
  algorithm/proto/algorithm_io_pb2_grpc.py