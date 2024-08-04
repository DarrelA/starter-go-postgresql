#!/bin/bash

COVDATAFILES_DIR="/root/testdata/covdatafiles"
OUTPUT_HTML="/root/testdata/output/it_coverage.html"

# Setup
rm "$OUTPUT_HTML"
rm -rf "$COVDATAFILES_DIR"
mkdir "$COVDATAFILES_DIR"

# Run test
GOCOVERDIR="$COVDATAFILES_DIR" ./integration_test.sh
ls "$COVDATAFILES_DIR" # @TODO: Figure out coverage profiling support for integration tests

# Check if coverage files are generated
if [ -z "$(ls -A $COVDATAFILES_DIR)" ]; then
    echo "No coverage data files generated."
    exit 1
fi

# Generate coverage data percentage in html
go tool covdata percent -i="$COVDATAFILES_DIR" -o "$COVDATAFILES_DIR/profile.txt"
go tool cover -html="$COVDATAFILES_DIR/profile.txt" -o "$OUTPUT_HTML"

echo "Completed running the wrap_test_for_coverage.sh"