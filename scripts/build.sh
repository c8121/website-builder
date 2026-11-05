#!/bin/bash

BASE_DIR=$(realpath "$(dirname "$0")")
OUT_DIR=$BASE_DIR/../bin
CMD_DIR=$BASE_DIR/../cmd

echo Build into $OUT_DIR
echo Sources from $CMD_DIR

mkdir -p $OUT_DIR

go build -o $OUT_DIR/spa-server $CMD_DIR/spa_server/main.go
go build -o $OUT_DIR/builder $CMD_DIR/builder/main.go
