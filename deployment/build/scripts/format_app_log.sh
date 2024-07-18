#!/bin/bash

# Read the input file and format the logs
input_file="deployment/logs/app.log"
output_file="deployment/logs/app.log.json"

# Wrap log entries in JSON array brackets and add commas between entries
{
  echo "["
  sed -e '$!s/$/,/' $input_file
  echo "]"
} > $output_file
