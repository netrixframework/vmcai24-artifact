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
docker build -f <benchmark>.dockerfile -t <benchmark>-netrix .
```

To run

``` bash
docker run <benchmark>-netrix:latest [ARGS]
```

The arguments passed are based on the benchmark that you are trying to reproduce.

### To run unit tests

To run the primary benchmarks of Table 3 from the paper for a specific unit test `UT`, please use the arguments `pct-test [UT]`
