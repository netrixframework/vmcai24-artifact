#!/bin/bash

CLIENTPATH=/netrixframework/java-clientlibrary
BFTSMARTPATH=/netrixframework/bft-smart

rm -rf ~/.m2/repository/io/github/netrixframework
cd $CLIENTPATH
gradle build
gradle publishToMavenLocal

cd $BFTSMARTPATH
rm config/currentView*
./gradlew clean
./gradlew installDist

mkdir -p $BFTSMARTPATH/build/logs