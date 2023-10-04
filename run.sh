#!/bin/bash

USAGE="Usage: run.sh <bftsmart/tendermint/raft> <test_case> <iterations>"

if [ $# -lt 2 ] || [ $# -gt 3 ]; then 
    echo $USAGE
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
TESTCASE=$2
ITERATIONS=100

# Checking if iteration is an integer
if [ $# -eq  3 ] && ! [[ $var =~ ^-?[0-9]+$ ]]; then 
    echo "error: invalid value for iterations argument"
    echo $USAGE
    exit 1
elif [ $# -eq  3 ]; then
    ITERATION=$3
fi

# Checking for valid test case strings

is_bftsmart_test() {
    if ! [ "$1" = "bftsmart" ]; then
        return 1
    else
        case "$2" in
            ("DPropForP") return 0 ;;
            ("DPropSame") return 0 ;;
            ("DropWrite") return 0 ;;
            ("DropWriteForP") return 0 ;;
            ("ExpectNewEpoch") return 0 ;;
            ("ExpectStop") return 0 ;;
            ("ByzLeaderChange") return 0 ;;
            ("PrevEpochProposal") return 0 ;;
            (*) return 1 ;;
        esac
    fi
}

is_tendermint_test() {
    if ! [ "$1" = "tendermint" ]; then
        return 1
    else
        case "$2" in
            ("ExpectUnlock") return 0 ;;
            ("Relocked") return 0 ;;
            ("LockedCommit") return 0 ;;
            ("LaggingReplica") return 0 ;;
            ("ForeverLaggingReplica") return 0 ;;
            ("RoundSkip") return 0 ;;
            ("BlockVotes") return 0 ;;
            ("PrecommitInvariant") return 0 ;;
            ("CommitAfterRoundSkip") return 0 ;;
            ("DifferentDecisions") return 0 ;;
            ("NilPrevotes") return 0 ;;
            ("ProposalNilPrevote") return 0 ;;
            ("NotNilDecide") return 0 ;;
            ("GarbledMessage") return 0 ;;
            ("HigherRound") return 0 ;;
            (*) return 1 ;;
        esac
    fi
}

is_raft_test() {
    if ! [ "$1" = "raft" ]; then
        return 1
    else
        case "$2" in
            ("Liveness") return 0 ;;
            ("LivenessNoCQ") return 0 ;;
            ("NoLiveness") return 0 ;;
            ("ConfChangeBug") return 0 ;;
            ("DropHeartbeat") return 0 ;;
            ("DropVotes") return 0 ;;
            ("DropFVotes") return 0 ;;
            ("DropAppend") return 0 ;;
            ("ReVote") return 0 ;;
            ("ManyReVote") return 0 ;;
            ("MultiReVote") return 0 ;;
            (*) return 1 ;;
        esac
    fi
}

is_valid_test() {
    case "$1" in
        ("tendermint") if is_tendermint_test $1 $2; then return 0; else return 1; fi ;;
        ("raft") if is_raft_test $1 $2; then return 0; else return 1; fi ;;
        ("bftsmart") if is_bftsmart_test $1 $2; then return 0; else return 1; fi ;;
        (*) return 1 ;;
    esac
}

if ! is_valid_test $BENCHMARK $TESTCASE; then
    echo "error: invalid test case arguments"
    echo $USAGE
    exit 1
fi


DOCKER_COMMAND=docker

if command -v podman &> /dev/null; then
    DOCKER_COMMAND=podman
fi

$DOCKER_COMMAND run -it "$DOCKER_IMAGE:latest" pct-test $TESTCASE -i $ITERATIONS


