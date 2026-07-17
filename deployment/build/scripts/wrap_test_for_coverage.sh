#!/usr/bin/env bash

set -Eeuo pipefail

readonly COVDATAFILES_DIR="/root/testdata/output/covdatafiles"
readonly OUTPUT_PROFILE="${COVDATAFILES_DIR}/coverage.out"

rm -rf "$COVDATAFILES_DIR"
mkdir -p "$COVDATAFILES_DIR"

GOCOVERDIR="$COVDATAFILES_DIR" ./integration_test.sh

go tool covdata percent -i="$COVDATAFILES_DIR"
go tool covdata textfmt -i="$COVDATAFILES_DIR" -o "$OUTPUT_PROFILE"

echo "Generated integration coverage profile at ${OUTPUT_PROFILE}"
