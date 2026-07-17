#!/usr/bin/env bash

set -Eeuo pipefail

readonly input_file="deployment/logs/app.log"
readonly output_file="deployment/logs/app.log.json"

mkdir -p "$(dirname "$output_file")"
if [[ ! -s "$input_file" ]]; then
    printf '[]\n' > "$output_file"
    exit 0
fi

{
    printf '[\n'
    sed -e '$!s/$/,/' "$input_file"
    printf ']\n'
} > "$output_file"
