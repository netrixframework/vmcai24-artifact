#!/usr/bin/bash

CODEHOME=/go/src/github.com/netrixframework/tendermint-testing
LOGPATH=/go/src/github.com/netrixframework/tendermint-testing/results

TMCODEHOME=/go/src/github.com/tendermint/tendermint

if [ $# -eq 1 ] && [ $1 = "--help" ]; then
    cd $CODEHOME
    ./tendermint-testing --help
    exit 1
fi

if [ ! -d "$CODEHOME" ]; then
    echo "Tendermint testing code does not exist"
    exit 1
fi

mkdir -p $LOGPATH
mkdir -p $LOGPATH/coverage
: > $LOGPATH/checker.log

echo "Starting Netrix server..."
cd $CODEHOME
./tendermint-testing $@ &
NETRIXPID=$!

tail -f $LOGPATH/checker.log | sed '/Waiting for all replicas to connect/ q'

echo "Starting tendermint servers..."
cd $TMCODEHOME
goreman start 2>&1 > /dev/null &

tail -f $LOGPATH/checker.log
wait $NETRIXPID
echo "Completed"
echo "Stopping tendermint..."
goreman run stop-all
echo "Done."