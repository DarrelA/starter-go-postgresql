#!/bin/bash

OUTPUT_JSON_FILE="./testdata/output/responses.json"

# Testdata files
TESTDATA_JSON_FILES=(
    "./testdata/json/register.json"
    "./testdata/json/login.json"
)

# Endpoints
URLS=(
    "localhost:8080/auth/api/v1/users/register"
    "localhost:8080/auth/api/v1/users/login"
)

# Run the main application
./starter-go-postgresql-it &

# Wait for the application to start
sleep 5

# Function to parse JSON using jq
parse_json() {
    echo "$1" | jq -r "$2"
}

# Initialize the output JSON file with an empty array
echo "[]" > "$OUTPUT_JSON_FILE"

# Create response_body.txt file
touch response_body.txt

# Loop through the JSON files and URLs
for i in "${!TESTDATA_JSON_FILES[@]}"; do
    TESTDATA_JSON_FILE=${TESTDATA_JSON_FILES[$i]}
    URL=${URLS[$i]}

    # Loop through the JSON file
    jq -c '.[]' "$TESTDATA_JSON_FILE" | while read -r item; do
        # Extract TestName, ExpectedStatusCode, and Input
        test_name=$(echo "$item" | jq -r '.TestName')
        expected_status_code=$(echo "$item" | jq -r '.ExpectedStatusCode')
        input=$(echo "$item" | jq -c '.Input')

        # Send POST request and capture response body and status
        response_status=$(curl -s \
            -o response_body.txt \
            -w "%{http_code}" \
            -X POST \
            -H "Content-Type: application/json" \
            -d "$input" \
            "$URL")

        # Read response body from the file
        response_body=$(cat response_body.txt)

        # Append the result to the output JSON file
        jq --arg test_name "$test_name" \
           --arg expected_status_code "$expected_status_code" \
           --arg response_status "$response_status" \
           --arg response_body "$response_body" \
           '. += [{"TestName": $test_name, "ExpectedStatusCode": $expected_status_code, "ResponseStatus": $response_status, "ResponseBody": $response_body}]' \
           "$OUTPUT_JSON_FILE" > temp.json && mv temp.json "$OUTPUT_JSON_FILE"

    done
done

# Cleanup temporary file
rm response_body.txt

echo "Completed running the integration_test.sh"