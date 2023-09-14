#!/usr/bin/bash

CODEHOME=/go/src/github.com/netrixframework/bftsmart-testing
LOGPATH=/go/src/github.com/netrixframework/bftsmart-testing/results

BFTSMARTHOME=/netrixframework/bft-smart

if [ $# -eq 1 ] && [ $1 = "--help" ]; then
    cd $CODEHOME
    ./bftsmart-testing --help
    exit 1
fi

if [ ! -d "$CODEHOME" ]; then
    echo "BFTSmart testing code does not exist"
    exit 1
fi

mkdir -p $LOGPATH
mkdir -p $LOGPATH/coverage
: > $LOGPATH/checker.log

echo "Starting Netrix server..."
cd $CODEHOME
./bftsmart-testing $@ &
NETRIXPID=$!

tail -f $LOGPATH/checker.log | sed '/Waiting for all replicas to connect/ q'

echo "Starting BFTSmart..."
cd $BFTSMARTHOME
goreman start 2>&1 > /dev/null &

tail -f $LOGPATH/checker.log
wait $NETRIXPID
echo "Completed"
echo "Stopping BFTSmart..."
goreman run stop-all
echo "Done."