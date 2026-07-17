#!/usr/bin/env bash

set -Eeuo pipefail

script_dir=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
repo_root=$(cd "${script_dir}/../../.." && pwd)
test_dir=$(mktemp -d)
trap 'rm -rf "$test_dir"' EXIT
export OUTPUT_JSON_FILE="${test_dir}/responses.json"
if [[ -d "${repo_root}/testdata/json" ]]; then
	export TESTDATA_DIR="${repo_root}/testdata/json"
else
	export TESTDATA_DIR="${script_dir}/testdata/json"
fi
# shellcheck source=integration_test.sh
source "${script_dir}/integration_test.sh"

status_matches 200 200
validate_testdata

if response_exposes_tokens '{"status":"success"}'; then
	echo "response_exposes_tokens rejected a token-free response" >&2
	exit 1
fi
if ! response_exposes_tokens '{"status":"success","access_token":"secret"}'; then
	echo "response_exposes_tokens accepted an access token" >&2
	exit 1
fi
if ! response_exposes_tokens '{"data":{"refresh_token":"secret"}}'; then
	echo "response_exposes_tokens accepted a nested refresh token" >&2
	exit 1
fi

if status_matches 200 500; then
    echo "status_matches accepted different status codes" >&2
    exit 1
fi

printf '[]\n' > "$OUTPUT_JSON_FILE"
record_result "/test" "passing case" 200 200 '{"status":"success"}'
record_result "/test" "failing case" 200 500 '{"status":"error"}'
jq --exit-status '
    length == 2 and
    .[0].Passed == true and
    .[0].ExpectedStatus == 200 and
    .[1].Passed == false and
    .[1].ActualStatus == 500
' "$OUTPUT_JSON_FILE" >/dev/null

bash -n "${script_dir}/integration_test.sh"
bash -n "${script_dir}/wrap_test_for_coverage.sh"

echo "Integration harness self-test passed"
