#!/bin/bash

if [ $# -ne 1 ]; then 
    echo "Usage: build.sh <bftsmart/tendermint/raft>"
    exit 1
fi

BENCHMARK=$1

is_valid_benchmark() {
    if [ "$1" = "tendermint" ] || [ "$1" = "raft" ] || [ "$1" = "bftsmart" ]; then
        return 0
    else 
        return 1
    fi
}

if ! is_valid_benchmark $BENCHMARK; then
    echo "error: invalid benchmark"
    echo $USAGE
    exit 1
fi

DOCKER_FILE="$BENCHMARK.dockerfile"
DOCKER_IMAGE="$BENCHMARK-netrix"

DOCKER_COMMAND=docker

if command -v podman &> /dev/null; then
    DOCKER_COMMAND=podman
fi

$DOCKER_COMMAND build -f $DOCKER_FILE -t $DOCKER_IMAGE .