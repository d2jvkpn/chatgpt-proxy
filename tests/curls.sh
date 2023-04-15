#! /usr/bin/env bash
set -eu -o pipefail
_wd=$(pwd)
_path=$(dirname $0 | xargs -i readlink -f {})

addr=http://localhost:3021

curl -i -X GET $addr/api/v1/open/version

curl -i -X POST $addr/v1/chat/completions \
  -d '{"messages": [{"role":"user", "content":"What are the key elements of the best automotive photography?"}]}'
