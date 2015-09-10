#!/bin/bash

cd ../Common
COMMONDIR=`pwd`
export GOPATH=$GOPATH':'$COMMONDIR
cd ../Client

CURDIR=`pwd`

export GOPATH=$GOPATH':'$COMMONDIR':'$CURDIR
echo $GOPATH

go build ./src/Client

rm -rf bin/log
mkdir -p bin/log

mv Client ./bin/gxclient




