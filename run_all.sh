#!/bin/bash

USAGE="Usage: run_all.sh <bftsmart/tendermint/raft>"

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

DOCKER_IMAGE="$BENCHMARK-netrix"
DOCKER_COMMAND=docker

if command -v podman &> /dev/null; then
    DOCKER_COMMAND=podman
fi

declare -a BFTSMART_TESTS=("DPropForP" "DPropSame" "DropWrite" "DropWriteForP" "ExpectNewEpoch" "ExpectStop" "ByzLeaderChange" "PrevEpochProposal")
declare -a TENDERMINT_TESTS=("ExpectUnlock" "Relocked" "LockedCommit" "LaggingReplica" "ForeverLaggingReplica" "RoundSkip" "BlockVotes" "PrecommitInvariant" "CommitAfterRoundSkip" "DifferentDecisions" "NilPrevotes" "ProposalNilPrevote" "NotNilDecide" "GarbledMessage" "HigherRound")
declare -a RAFT_TESTS=("Liveness" "LivenessNoCQ" "NoLiveness" "ConfChangeBug" "DropHeartbeat" "DropVotes" "DropFVotes" "DropAppend" "ReVote" "ManyReVote" "MultiReVote")

TESTS=$BFTSMART_TESTS
case "$BENCHMARK" in
    ("tendermint") TESTS=$TENDERMINT_TESTS ;;
    ("bftsmart") TESTS=$BFTSMART_TESTS ;;
    ("raft") TESTS=$RAFT_TESTS ;;
esac

ITERATIONS=100

for i in "${TESTS[@]}"
do
    echo "Running test: $i"
    $DOCKER_COMMAND run -it $DOCKER_IMAGE:latest pct-test $i -i $ITERATIONS
done