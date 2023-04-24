#! /usr/bin/env bash
set -eu -o pipefail
_wd=$(pwd)
_path=$(dirname $0 | xargs -i readlink -f {})

image=registry.cn-shanghai.aliyuncs.com/d2jvkpn/chatgpt-proxy-py
tag=latest

docker build -f ${_path}/Dockerfile-py --tag $image:$tag ./

docker push $image:$tag
