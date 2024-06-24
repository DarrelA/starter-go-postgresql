#!/bin/bash

# Generate 2048-bit RSA Private Key for Access Token
openssl genpkey -algorithm RSA -out access_private_key.pem -pkeyopt rsa_keygen_bits:2048

# Extract the Public Key from the Private Key for Access Token
openssl rsa -pubout -in access_private_key.pem -out access_public_key.pem

# Generate 2048-bit RSA Private Key for Refresh Token
openssl genpkey -algorithm RSA -out refresh_private_key.pem -pkeyopt rsa_keygen_bits:2048

# Extract the Public Key from the Private Key for Refresh Token
openssl rsa -pubout -in refresh_private_key.pem -out refresh_public_key.pem

# Encode the Private and Public Keys to Base64 (using cat to ensure compatibility)
access_private_key_base64=$(cat access_private_key.pem | base64 | tr -d '\n')
access_public_key_base64=$(cat access_public_key.pem | base64 | tr -d '\n')
refresh_private_key_base64=$(cat refresh_private_key.pem | base64 | tr -d '\n')
refresh_public_key_base64=$(cat refresh_public_key.pem | base64 | tr -d '\n')

# Output the keys to a .txt file in the desired format
cat <<EOL > refresh_token_keys.txt
ACCESS_TOKEN_PRIVATE_KEY=$access_private_key_base64

ACCESS_TOKEN_PUBLIC_KEY=$access_public_key_base64

REFRESH_TOKEN_PRIVATE_KEY=$refresh_private_key_base64

REFRESH_TOKEN_PUBLIC_KEY=$refresh_public_key_base64
EOL

# Clean up the generated key files
rm access_private_key.pem access_public_key.pem refresh_private_key.pem refresh_public_key.pem

echo "Keys generated and saved to refresh_token_keys.txt"