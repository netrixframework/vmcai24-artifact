# vmcai24-artifact

Artifact for VMCAI24 paper "A Domain Specific Language for Testing Distributed Protocol Implementations"

The corresponding github [repo](https://github.com/netrixframework/vmcai24-artifact)

## Organization

The folders associated with each of the benchmarks are:

- BFTSmart - instrumented code in `bft-smart` and tests in`bftsmart-testing`
- Tendermint - instrumented code in `tendermint` and tests in `tendermint-testing`
- Raft - instrumented code `raft-testing/raft` and the tests in `raft-testing/tests`

In addition to the code and scripts to reproduce the tests, we include the run logs of the test cases for in `run_logs` and segregate it based on the benchmark.

The `scripts` folder contains the scripts used to run the tests.

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

The log prints out the number of successful outcomes after every iteration.

To run all the tests for a given benchmark

```bash
./run_all.sh <benchmark>
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

## Experimental setup

Each iteration of the test varies per the benchmarks, based on the setup time. This also gives an estimate for the run time of the complete test.

- For Etcd-raft, the iteration timeout is 60 seconds.
- For Tendermint, 90 seconds.
- For BFTSmart, 60 seconds.

The timeout can be changed in `<benchmark>-testing/cmd/pct_test_strat.go` for the key `IterationTimeout`.

The build script builds the docker image. However, the performance of the tests is sensitive to the resources available to the docker container. In each test iteration, we run the implementation and communicate events/messages to the central server. The number of communication messages exchanged is large and in the order of 10k. Therefore running the docker image with minimal bandwidth or memory resources increases the message exchange time. The test timeout should be increased as a consequence.

We expect some hiccups with the BFT-Smart benchmark as it is sensitive to the memory requirements. The test might crash without complete all the iterations due to inadequate memory.


## Extending tests

For each benchmark, `<benchmark>-testing` directory contains the set of tests. More tests for the benchmark can be added in the `<benchmark>-testing/tests` directory.

Corresponding to each new test, the `<benchmark>-testing/cmd/strategies.go` needs to be updated with the new test name, filters and the state machine.

Additionally, while running the tests, `<benchmark>-testing/cmd/pct_test_strat.go` contains the parameters that can be modified.
