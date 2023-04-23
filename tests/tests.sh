#! /usr/bin/env bash
set -eu -o pipefail
_wd=$(pwd)
_path=$(dirname $0 | xargs -i readlink -f {})

#### tls
mkdir -p configs

openssl genrsa -out configs/server.key 2048

openssl req -new \
  -subj "/CN=127.0.0.1/O=My Organization/C=US" \
  -key configs/server.key \
  -out configs/server.csr

openssl x509 -req -days 365 \
  -in configs/server.csr \
  -signkey configs/server.key \
  -out configs/server.crt

# using server.key and server.crt for TLS encryption

#### generate api-keys, sk-x32
tr -dc A-Za-z0-9 </dev/urandom | head -c 32; echo ''

#### service
go run main.go &

#### test
# curl -i ... skip tls verification
addr=http://localhost:3021
api_key="sk-xxxxx"

time curl -i -X GET $addr/api/v1/open/version

time curl -i -X POST $addr/v1/chat/completions \
  -H "Authorization: Bearer $api_key" \
  -d '{"model": "gpt-3.5-turbo", "messages": [{"role":"user", "content":"Who are you?"}]}'
