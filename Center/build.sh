#!/bin/bash

cd ../Common
COMMONDIR=`pwd`
export GOPATH=$GOPATH':'$COMMONDIR
cd ../Center

CURDIR=`pwd`

export GOPATH=$GOPATH':'$COMMONDIR':'$CURDIR
echo $GOPATH

go build ./src/Center

rm -rf bin/log
mkdir -p bin/log

mv Center ./bin/gxcenter