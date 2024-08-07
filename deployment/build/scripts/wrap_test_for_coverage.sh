#!/bin/bash

COVDATAFILES_DIR="/root/testdata/output/covdatafiles"
OUTPUT_HTML="/root/testdata/output/it_coverage.html"

# Setup
rm "$OUTPUT_HTML"
rm -rf "$COVDATAFILES_DIR"
mkdir "$COVDATAFILES_DIR"

# Run test
GOCOVERDIR="$COVDATAFILES_DIR" ./integration_test.sh

# Generate coverage data percentage in html
go tool covdata percent -i="$COVDATAFILES_DIR" -o "$COVDATAFILES_DIR/coverage.out"

echo "Completed running the wrap_test_for_coverage.sh"