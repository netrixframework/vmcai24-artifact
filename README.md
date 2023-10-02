# vmcai24-artifact

Artifact for VMCAI24 paper "A Domain Specific Language for Testing Distributed Protocol Implementations"

## Organization

The folders associated with each of the benchmarks are:

- BFTSmart - instrumented code in `bft-smart` and tests in`bftsmart-testing`
- Tendermint - instrumented code in `tendermint` and tests in `tendermint-testing`
- Raft - instrumented code `raft-testing/raft` and the tests in `raft-testing/tests`

Additionally,

- `scripts` contain build and run scripts used within the dockerfile for each of the benchmarks.
- `netrix` contains the framework code. Please refer [here](https://netrixframework.github.io) for documentation
- `java-clientlibrary` contains code for the java client used by the instrumented BFTSmart library code.

## Setup and run

For each benchmark you can build the dockerfile and run it with the following commands

To build

``` bash
./build.sh <tendermint/raft/bftsmart>
```

To run for a default of 100 iterations

``` bash
./run.sh <tendermint/raft/bftsmart> <test_case> 
```

Optionally, you can pass the number of iterations as the third argument to the run script

``` bash
./run.sh <tendermint/raft/bftsmart> <test_case> <iterations>
```

### To run unit tests

The run script allows reproducing the expected results for the tests based on the benchmarks listed in table 3 of the paper. We document below the list of unit tests for each benchmark that can be passed as arguments

For **raft**,

- Liveness
- LivenessNoCQ
- NoLiveness
- ConfChangeBug
- DropHeartbeat
- DropVotes
- DropFVotes
- DropAppend
- ReVote
- ManyReVote
- MultiReVote

For **tendermint**,

- ExpectUnlock
- Relocked
- LockedCommit
- LaggingReplica
- ForeverLaggingReplica
- RoundSkip
- BlockVotes
- PrecommitInvariant
- CommitAfterRoundSkip
- DifferentDecisions
- NilPrevotes
- ProposalNilPrevote
- NotNilDecide
- GarbledMessage
- HigherRound

For **bftsmart**,

- DPropForP
- DPropSame
- DropWrite
- DropWriteForP
- ExpectNewEpoch
- ExpectStop
- ByzLeaderChange
- PrevEpochProposal
