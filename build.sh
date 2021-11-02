#!/bin/bash
set -xeuo pipefail

GOOS=js GOARCH=wasm go build -o main.wasm ./main.go
