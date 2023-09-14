#!/usr/bin/bash

CODEHOME=/go/src/github.com/netrixframework/raft-testing
LOGPATH=/go/src/github.com/netrixframework/raft-testing/results

if [ $# -eq 1 ] && [ $1 = "--help" ]; then
    cd $CODEHOME
    ./raft-testing --help
    exit 1
fi

if [ ! -d "$CODEHOME" ]; then
    echo "Raft testing code does not exist"
    exit 1
fi

mkdir -p $LOGPATH
mkdir -p $LOGPATH/coverage
: > $LOGPATH/checker.log

echo "Starting Netrix server..."
cd $CODEHOME
./raft-testing $@ &
NETRIXPID=$!

tail -f $LOGPATH/checker.log | sed '/Waiting for all replicas to connect/ q'

echo "Starting raft..."
cd $CODEHOME/raft
goreman start 2>&1 > /dev/null &

tail -f $LOGPATH/checker.log
wait $NETRIXPID
echo "Completed"
echo "Stopping raft..."
goreman run stop-all
echo "Done."