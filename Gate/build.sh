#!/bin/bash

cd ../Common
COMMONDIR=`pwd`
export GOPATH=$GOPATH':'$COMMONDIR
cd ../Gate

CURDIR=`pwd`

export GOPATH=$GOPATH':'$COMMONDIR':'$CURDIR
echo $GOPATH

go build ./src/Gate

rm -rf bin/log
mkdir -p bin/log

mv Gate ./bin/gxgate