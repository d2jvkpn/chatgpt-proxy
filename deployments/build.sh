#! /usr/bin/env bash
set -eu -o pipefail
_wd=$(pwd)
_path=$(dirname $0 | xargs -i readlink -f {})


function onExit {
    git checkout dev # --force
}
trap onExit EXIT

####
gitBranch=$1
image="registry.cn-shanghai.aliyuncs.com/d2jvkpn/chatgpt-proxy"
tag=$gitBranch

git checkout $gitBranch
git pull --no-edit

buildTime=$(date +'%FT%T%:z')
gitBranch="$(git rev-parse --abbrev-ref HEAD)" # current branch
gitCommit=$(git rev-parse --verify HEAD) # git log --pretty=format:'%h' -n 1
gitTime=$(git log -1 --format="%at" | xargs -I{} date -d @{} +%FT%T%:z)
gitTreeState="clean"

uncommitted=$(git status --short)
unpushed=$(git diff origin/$gitBranch..HEAD --name-status)
[[ ! -z "$uncommitted$unpushed" ]] && gitTreeState="dirty"

####
for base in $(awk '/^FROM/{print $2}' ${_path}/Dockerfile); do
    echo ">>> pull $bae"
    docker pull $base
    bn=$(echo $base | awk -F ":" '{print $1}')
    if [[ -z "$bn" ]]; then continue; fi
    docker images --filter "dangling=true" --quiet "$bn" | xargs -i docker rmi {}
done &> /dev/null

echo ">>> build image: $image:$tag..."

ldflags="-X main.buildTime=$buildTime -X main.gitBranch=$gitBranch \
  -X main.gitCommit=$gitCommit -X main.gitTime=$gitTime \
  -X main.gitTreeState=$gitTreeState"

df=${_path}/Dockerfile

docker build --no-cache --file $df --build-arg=ldflags="$ldflags" --tag $image:$tag ./

docker image prune --force --filter label=stage=chatgpt-proxy_builder &> /dev/null

#### push image
echo ">>> push image: $image:$tag..."
docker push $image:$tag

images=$(docker images --filter "dangling=true" --quiet $image)
for img in $images; do docker rmi $img || true; done &> /dev/null
