#! /usr/bin/env bash
set -eu -o pipefail
_wd=$(pwd)
_path=$(dirname $0 | xargs -i readlink -f {})

#### tls
mkdir -p configs

openssl genrsa -out configs/server.key 2048

openssl req -new \
  -subj "/CN=example.com/O=My Organization/C=US" \
  -key configs/server.key \
  -out configs/server.csr

openssl x509 -req -days 365 \
  -in configs/server.csr \
  -signkey configs/server.key \
  -out configs/server.crt

# using server.key and server.crt for TLS encryption

#### service
go run main.go &

#### test
addr=http://localhost:3021
# curl -i ... skip tls verification

curl -i -X GET $addr/api/v1/open/version

curl -i -X POST $addr/v1/chat/completions \
  -H "Authorization: Bearer sk-xxxxxxxx" \
  -d '{"messages": [{"role":"user", "content":"What are the key elements of the best automotive photography?"}]}'
