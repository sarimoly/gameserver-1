#!/bin/bash

cd ../Common
COMMONDIR=`pwd`
export GOPATH=$GOPATH':'$COMMONDIR
cd ../Tool

CURDIR=`pwd`

export GOPATH=$GOPATH':'$COMMONDIR':'$CURDIR
echo $GOPATH

go build ./src/NewServer
go build ./src/Test

rm -rf bin/log
mkdir -p bin/log

mv NewServer ./bin/gxNewServer
mv Test ./bin/gxTest