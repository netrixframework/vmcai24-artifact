#!/bin/bash

CLIENTPATH=/netrixframework/java-clientlibrary
BFTSMARTPATH=/netrixframework/bft-smart

rm -rf ~/.m2/repository/io/github/netrixframework
cd $CLIENTPATH
./gradlew --refresh-dependencies build
./gradlew publishToMavenLocal

cd $BFTSMARTPATH
rm config/currentView*
./gradlew clean
./gradlew --refresh-dependencies installDist

mkdir -p $BFTSMARTPATH/build/logs