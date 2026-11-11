#!/usr/bin/env zsh

if [[ $# -ne 2 ]]; then
  echo "Usage: $0 <url> <r2-path>"
  exit 1
fi

URL="$1"
R2_PATH="$2"

go run ./cmd/main.go \
  -r2-acckey="XXX" \
  -r2-bucket="XXX" \
  -r2-endpoint="XXX" \
  -r2-secretkey="XXX" \
  -r2-path="$R2_PATH" \
  -url="$URL"
