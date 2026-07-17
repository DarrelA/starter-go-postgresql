#!/usr/bin/env bash

set -euo pipefail
umask 077

output_file=${1:-refresh_token_keys.txt}
key_dir=$(mktemp -d)
trap 'rm -rf "$key_dir"' EXIT

generate_key_pair() {
	local name=$1
	openssl genpkey -algorithm RSA \
		-out "$key_dir/${name}_private.pem" \
		-pkeyopt rsa_keygen_bits:2048 >/dev/null 2>&1
	openssl rsa -pubout \
		-in "$key_dir/${name}_private.pem" \
		-out "$key_dir/${name}_public.pem" >/dev/null 2>&1
}

encode_file() {
	base64 < "$1" | tr -d '\r\n'
}

replace_env_value() {
	local key=$1
	local value=$2
	local file=$3
	local temporary_file
	temporary_file=$(mktemp)

	awk -v key="$key" -v value="$value" '
		BEGIN { found = 0 }
		$0 ~ ("^" key "=") { print key "=" value; found = 1; next }
		{ print }
		END { if (!found) print key "=" value }
	' "$file" > "$temporary_file"

	mv "$temporary_file" "$file"
}

generate_key_pair access
generate_key_pair refresh

access_private_key=$(encode_file "$key_dir/access_private.pem")
access_public_key=$(encode_file "$key_dir/access_public.pem")
refresh_private_key=$(encode_file "$key_dir/refresh_private.pem")
refresh_public_key=$(encode_file "$key_dir/refresh_public.pem")
csrf_secret=$(openssl rand -base64 32 | tr -d '\r\n')

if [[ -f "$output_file" ]]; then
	replace_env_value ACCESS_TOKEN_PRIVATE_KEY "$access_private_key" "$output_file"
	replace_env_value ACCESS_TOKEN_PUBLIC_KEY "$access_public_key" "$output_file"
	replace_env_value REFRESH_TOKEN_PRIVATE_KEY "$refresh_private_key" "$output_file"
	replace_env_value REFRESH_TOKEN_PUBLIC_KEY "$refresh_public_key" "$output_file"
	replace_env_value CSRF_SECRET "$csrf_secret" "$output_file"
else
	{
		printf 'ACCESS_TOKEN_PRIVATE_KEY=%s\n' "$access_private_key"
		printf 'ACCESS_TOKEN_PUBLIC_KEY=%s\n' "$access_public_key"
		printf 'REFRESH_TOKEN_PRIVATE_KEY=%s\n' "$refresh_private_key"
		printf 'REFRESH_TOKEN_PUBLIC_KEY=%s\n' "$refresh_public_key"
		printf 'CSRF_SECRET=%s\n' "$csrf_secret"
	} > "$output_file"
fi

echo "Generated JWT signing keys and CSRF secret in $output_file"
