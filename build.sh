#!/bin/bash

set -exo pipefail

cd "$(dirname "$0")" || exit 1

export DOCKER_BUILDKIT=1

image=${IMAGE:-chess/game}

build_cmd="docker build . --file Dockerfile --tag ${image} --build-arg BUILDKIT_INLINE_CACHE=1"

if [ -n "$CACHE_FROM" ]; then
    build_cmd+=" ${CACHE_FROM}"
fi

eval "$build_cmd"
