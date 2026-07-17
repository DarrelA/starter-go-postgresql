#!/usr/bin/env bash

set -Eeuo pipefail

readonly BASE_URL="${BASE_URL:-http://localhost:8080}"
readonly OUTPUT_JSON_FILE="${OUTPUT_JSON_FILE:-/root/testdata/output/responses.json}"
readonly HEALTH_URL="${BASE_URL}/auth/health"
readonly STARTUP_ATTEMPTS="${STARTUP_ATTEMPTS:-30}"
readonly TESTDATA_DIR="${TESTDATA_DIR:-/root/testdata/json}"

readonly -a TESTDATA_JSON_FILES=(
    "${TESTDATA_DIR}/register.json"
    "${TESTDATA_DIR}/login.json"
)

readonly -a ENDPOINTS=(
    "/auth/api/v1/users/register"
    "/auth/api/v1/users/login"
)

app_pid=""
response_body_file=""
cookie_jar_file=""

cleanup() {
	if [[ -n "$response_body_file" ]]; then
		rm -f "$response_body_file"
	fi
	if [[ -n "$cookie_jar_file" ]]; then
		rm -f "$cookie_jar_file"
	fi
    rm -f "${OUTPUT_JSON_FILE}.tmp"
    if [[ -n "$app_pid" ]] && kill -0 "$app_pid" 2>/dev/null; then
        kill -TERM "$app_pid" 2>/dev/null || true
        wait "$app_pid" 2>/dev/null || true
    fi
}

status_matches() {
	[[ "$1" == "$2" ]]
}

response_exposes_tokens() {
	jq --exit-status '
		[.. | objects | select(has("access_token") or has("refresh_token"))] | length > 0
	' <<<"$1" >/dev/null 2>&1
}

cookie_jar_contains() {
	local cookie_name=$1
	awk -v name="$cookie_name" '$6 == name { found = 1 } END { exit !found }' "$cookie_jar_file"
}

cookie_value() {
	local jar_file=$1 cookie_name=$2
	awk -v name="$cookie_name" '$6 == name { value = $7 } END { print value }' "$jar_file"
}

validate_testdata() {
    local file
    if ((${#TESTDATA_JSON_FILES[@]} != ${#ENDPOINTS[@]})); then
        echo "Each test-data file must have a matching endpoint" >&2
        return 1
    fi
    for file in "${TESTDATA_JSON_FILES[@]}"; do
        jq --exit-status '
            type == "array" and
            length > 0 and
            all(.[];
                (.TestName | type == "string" and length > 0) and
                (.ExpectedStatusCode | type == "number") and
                (.Input | type == "object")
            )
        ' "$file" >/dev/null
    done
}

wait_for_application() {
    local attempt
    for ((attempt = 1; attempt <= STARTUP_ATTEMPTS; attempt++)); do
        if curl --silent --fail --max-time 2 "$HEALTH_URL" >/dev/null; then
            return 0
        fi
        if ! kill -0 "$app_pid" 2>/dev/null; then
            wait "$app_pid" || true
            echo "Application exited before becoming ready" >&2
            return 1
        fi
        sleep 1
    done
    echo "Application did not become ready after ${STARTUP_ATTEMPTS} attempts" >&2
    return 1
}

record_result() {
    local endpoint=$1
    local test_name=$2
    local expected_status=$3
    local actual_status=$4
    local response_body=$5

    jq \
        --arg endpoint "$endpoint" \
        --arg test_name "$test_name" \
        --arg expected_status "$expected_status" \
        --arg actual_status "$actual_status" \
        --arg response_body "$response_body" \
        '. += [{
            "Endpoint": $endpoint,
            "TestName": $test_name,
            "ExpectedStatus": ($expected_status | tonumber),
            "ActualStatus": ($actual_status | tonumber),
            "Passed": ($expected_status == $actual_status),
            "ResponseBody": $response_body
        }]' \
        "$OUTPUT_JSON_FILE" > "${OUTPUT_JSON_FILE}.tmp"
    mv "${OUTPUT_JSON_FILE}.tmp" "$OUTPUT_JSON_FILE"
}

run_cases() {
    local failures=0
    local index item test_name expected_status input actual_status response_body endpoint

    for index in "${!TESTDATA_JSON_FILES[@]}"; do
        endpoint="${ENDPOINTS[$index]}"
		while IFS= read -r item; do
            test_name=$(jq -r '.TestName' <<<"$item")
            expected_status=$(jq -r '.ExpectedStatusCode' <<<"$item")
            input=$(jq -c '.Input' <<<"$item")

            : > "$response_body_file"
			if ! actual_status=$(curl --silent --show-error \
				--output "$response_body_file" \
				--write-out "%{http_code}" \
				--max-time 10 \
				--cookie "$cookie_jar_file" \
				--cookie-jar "$cookie_jar_file" \
				--request POST \
                --header "Content-Type: application/json" \
                --data "$input" \
                "${BASE_URL}${endpoint}"); then
                actual_status="000"
            fi
            response_body=$(<"$response_body_file")
            record_result "$endpoint" "$test_name" "$expected_status" "$actual_status" "$response_body"

			if ! status_matches "$expected_status" "$actual_status"; then
				echo "FAIL: ${test_name}: expected HTTP ${expected_status}, got ${actual_status}" >&2
				failures=$((failures + 1))
			fi
			if response_exposes_tokens "$response_body"; then
				echo "FAIL: ${test_name}: response exposed an authentication token" >&2
				failures=$((failures + 1))
			fi
		done < <(jq -c '.[]' "${TESTDATA_JSON_FILES[$index]}")
    done

    if ((failures > 0)); then
        echo "${failures} integration case(s) failed; see ${OUTPUT_JSON_FILE}" >&2
        return 1
	fi
}

verify_current_user() {
	local login_input login_status me_status response_body
	login_input=$(jq -c '.[0].Input' "${TESTDATA_JSON_FILES[1]}")
	: > "$cookie_jar_file"
	login_status=$(curl --silent --show-error \
		--output "$response_body_file" \
		--write-out "%{http_code}" \
		--max-time 10 \
		--cookie-jar "$cookie_jar_file" \
		--request POST \
		--header "Content-Type: application/json" \
		--data "$login_input" \
		"${BASE_URL}${ENDPOINTS[1]}")
	if ! status_matches 200 "$login_status"; then
		echo "FAIL: current-user setup login: expected HTTP 200, got ${login_status}" >&2
		return 1
	fi

	: > "$response_body_file"
	me_status=$(curl --silent --show-error \
		--output "$response_body_file" \
		--write-out "%{http_code}" \
		--max-time 10 \
		--cookie "$cookie_jar_file" \
		--request GET \
		"${BASE_URL}/auth/api/v1/users/me")
	response_body=$(<"$response_body_file")
	record_result "/auth/api/v1/users/me" "valid session loads current user" 200 "$me_status" "$response_body"
	if ! status_matches 200 "$me_status"; then
		echo "FAIL: current user: expected HTTP 200, got ${me_status}" >&2
		return 1
	fi
	if ! jq --exit-status '
		.data.userRecord.uuid | type == "string" and length > 0
	' <<<"$response_body" >/dev/null; then
		echo "FAIL: current-user response did not contain a UUID" >&2
		return 1
	fi
	if response_exposes_tokens "$response_body"; then
		echo "FAIL: current-user response exposed an authentication token" >&2
		return 1
	fi
}

verify_logout_contract() {
	local csrf_token login_input login_status get_status post_status response_body
	login_input=$(jq -c '.[0].Input' "${TESTDATA_JSON_FILES[1]}")
	: > "$response_body_file"
	login_status=$(curl --silent --show-error \
		--output "$response_body_file" \
		--write-out "%{http_code}" \
		--max-time 10 \
		--cookie "$cookie_jar_file" \
		--cookie-jar "$cookie_jar_file" \
		--request POST \
		--header "Content-Type: application/json" \
		--data "$login_input" \
		"${BASE_URL}${ENDPOINTS[1]}")
	response_body=$(<"$response_body_file")
	if ! status_matches 200 "$login_status"; then
		echo "FAIL: logout setup login: expected HTTP 200, got ${login_status}" >&2
		return 1
	fi
	if response_exposes_tokens "$response_body"; then
		echo "FAIL: logout setup login exposed an authentication token" >&2
		return 1
	fi
	for cookie_name in access_token refresh_token csrf_token; do
		if ! cookie_jar_contains "$cookie_name"; then
			echo "FAIL: logout setup login did not issue ${cookie_name} cookie" >&2
			return 1
		fi
	done
	csrf_token=$(cookie_value "$cookie_jar_file" csrf_token)

	: > "$response_body_file"
	get_status=$(curl --silent --show-error \
		--output "$response_body_file" \
		--write-out "%{http_code}" \
		--max-time 10 \
		--cookie "$cookie_jar_file" \
		--request GET \
		"${BASE_URL}/auth/api/v1/users/logout")
	response_body=$(<"$response_body_file")
	record_result "/auth/api/v1/users/logout" "GET logout is not routed" 404 "$get_status" "$response_body"
	if ! status_matches 404 "$get_status"; then
		echo "FAIL: GET logout: expected HTTP 404, got ${get_status}" >&2
		return 1
	fi

	: > "$response_body_file"
	post_status=$(curl --silent --show-error \
		--output "$response_body_file" \
		--write-out "%{http_code}" \
		--max-time 10 \
		--cookie "$cookie_jar_file" \
		--cookie-jar "$cookie_jar_file" \
		--request POST \
		--header "X-CSRF-Token: ${csrf_token}" \
		"${BASE_URL}/auth/api/v1/users/logout")
	response_body=$(<"$response_body_file")
	record_result "/auth/api/v1/users/logout" "POST logout revokes session" 200 "$post_status" "$response_body"
	if ! status_matches 200 "$post_status"; then
		echo "FAIL: POST logout: expected HTTP 200, got ${post_status}" >&2
		return 1
	fi
	if response_exposes_tokens "$response_body"; then
		echo "FAIL: POST logout response exposed an authentication token" >&2
		return 1
	fi
}

verify_refresh_security() {
	local concurrent_dir csrf_token current_refresh login_input login_status new_refresh
	local no_csrf_status replay_status refresh_status revoked_status response_body
	local first_pid second_pid first_status second_status sorted_statuses
	local old_cookie_jar first_cookie_jar second_cookie_jar

	login_input=$(jq -c '.[0].Input' "${TESTDATA_JSON_FILES[1]}")
	: > "$cookie_jar_file"
	login_status=$(curl --silent --show-error --output "$response_body_file" --write-out "%{http_code}" \
		--max-time 10 --cookie-jar "$cookie_jar_file" --request POST \
		--header "Content-Type: application/json" --data "$login_input" \
		"${BASE_URL}${ENDPOINTS[1]}")
	if ! status_matches 200 "$login_status"; then
		echo "FAIL: refresh setup login: expected HTTP 200, got ${login_status}" >&2
		return 1
	fi
	csrf_token=$(cookie_value "$cookie_jar_file" csrf_token)
	current_refresh=$(cookie_value "$cookie_jar_file" refresh_token)
	if [[ -z "$csrf_token" || -z "$current_refresh" ]]; then
		echo "FAIL: refresh setup did not issue refresh and CSRF cookies" >&2
		return 1
	fi

	no_csrf_status=$(curl --silent --show-error --output "$response_body_file" --write-out "%{http_code}" \
		--max-time 10 --cookie "$cookie_jar_file" --request POST \
		"${BASE_URL}/auth/api/v1/users/refresh")
	if ! status_matches 403 "$no_csrf_status"; then
		echo "FAIL: refresh without CSRF: expected HTTP 403, got ${no_csrf_status}" >&2
		return 1
	fi

	old_cookie_jar=$(mktemp)
	cp "$cookie_jar_file" "$old_cookie_jar"
	refresh_status=$(curl --silent --show-error --output "$response_body_file" --write-out "%{http_code}" \
		--max-time 10 --cookie "$cookie_jar_file" --cookie-jar "$cookie_jar_file" --request POST \
		--header "X-CSRF-Token: ${csrf_token}" \
		"${BASE_URL}/auth/api/v1/users/refresh")
	if ! status_matches 200 "$refresh_status"; then
		echo "FAIL: refresh rotation: expected HTTP 200, got ${refresh_status}" >&2
		rm -f "$old_cookie_jar"
		return 1
	fi
	new_refresh=$(cookie_value "$cookie_jar_file" refresh_token)
	if [[ -z "$new_refresh" || "$new_refresh" == "$current_refresh" ]]; then
		echo "FAIL: refresh rotation did not replace the refresh token" >&2
		rm -f "$old_cookie_jar"
		return 1
	fi

	replay_status=$(curl --silent --show-error --output "$response_body_file" --write-out "%{http_code}" \
		--max-time 10 --cookie "$old_cookie_jar" --request POST \
		--header "X-CSRF-Token: ${csrf_token}" \
		"${BASE_URL}/auth/api/v1/users/refresh")
	rm -f "$old_cookie_jar"
	if ! status_matches 401 "$replay_status"; then
		echo "FAIL: refresh replay: expected HTTP 401, got ${replay_status}" >&2
		return 1
	fi
	csrf_token=$(cookie_value "$cookie_jar_file" csrf_token)
	revoked_status=$(curl --silent --show-error --output "$response_body_file" --write-out "%{http_code}" \
		--max-time 10 --cookie "$cookie_jar_file" --request POST \
		--header "X-CSRF-Token: ${csrf_token}" \
		"${BASE_URL}/auth/api/v1/users/refresh")
	if ! status_matches 401 "$revoked_status"; then
		echo "FAIL: replay did not revoke the replacement session; got HTTP ${revoked_status}" >&2
		return 1
	fi

	: > "$cookie_jar_file"
	login_status=$(curl --silent --show-error --output "$response_body_file" --write-out "%{http_code}" \
		--max-time 10 --cookie-jar "$cookie_jar_file" --request POST \
		--header "Content-Type: application/json" --data "$login_input" \
		"${BASE_URL}${ENDPOINTS[1]}")
	if ! status_matches 200 "$login_status"; then
		echo "FAIL: concurrent refresh setup login: expected HTTP 200, got ${login_status}" >&2
		return 1
	fi
	csrf_token=$(cookie_value "$cookie_jar_file" csrf_token)
	concurrent_dir=$(mktemp -d)
	first_cookie_jar="${concurrent_dir}/first.cookies"
	second_cookie_jar="${concurrent_dir}/second.cookies"
	cp "$cookie_jar_file" "$first_cookie_jar"
	cp "$cookie_jar_file" "$second_cookie_jar"
	curl --silent --show-error --output "${concurrent_dir}/first.body" --write-out "%{http_code}" \
		--max-time 10 --cookie "$first_cookie_jar" --cookie-jar "$first_cookie_jar" --request POST \
		--header "X-CSRF-Token: ${csrf_token}" \
		"${BASE_URL}/auth/api/v1/users/refresh" > "${concurrent_dir}/first.status" &
	first_pid=$!
	curl --silent --show-error --output "${concurrent_dir}/second.body" --write-out "%{http_code}" \
		--max-time 10 --cookie "$second_cookie_jar" --cookie-jar "$second_cookie_jar" --request POST \
		--header "X-CSRF-Token: ${csrf_token}" \
		"${BASE_URL}/auth/api/v1/users/refresh" > "${concurrent_dir}/second.status" &
	second_pid=$!
	wait "$first_pid"
	wait "$second_pid"
	first_status=$(<"${concurrent_dir}/first.status")
	second_status=$(<"${concurrent_dir}/second.status")
	sorted_statuses=$(printf '%s\n%s\n' "$first_status" "$second_status" | sort | tr '\n' ' ')
	rm -rf "$concurrent_dir"
	if [[ "$sorted_statuses" != "200 401 " ]]; then
		echo "FAIL: concurrent refresh expected one HTTP 200 and one HTTP 401, got ${first_status} and ${second_status}" >&2
		return 1
	fi
}

main() {
    mkdir -p "$(dirname "$OUTPUT_JSON_FILE")"
	response_body_file=$(mktemp)
	cookie_jar_file=$(mktemp)
	trap cleanup EXIT INT TERM
    printf '[]\n' > "$OUTPUT_JSON_FILE"
    validate_testdata

    ./starter-go-postgresql-it &
    app_pid=$!
	wait_for_application
	run_cases
	verify_current_user
	verify_refresh_security
	verify_logout_contract

    kill -TERM "$app_pid"
    wait "$app_pid"
    app_pid=""
    echo "Integration tests passed"
}

if [[ "${BASH_SOURCE[0]}" == "$0" ]]; then
    main "$@"
fi
