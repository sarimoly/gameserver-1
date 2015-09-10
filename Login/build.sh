#!/bin/bash

cd ../Common
COMMONDIR=`pwd`
export GOPATH=$GOPATH':'$COMMONDIR
cd ../Login

CURDIR=`pwd`

export GOPATH=$GOPATH':'$COMMONDIR':'$CURDIR
echo $GOPATH

go build ./src/Login

rm -rf bin/log
mkdir -p bin/log

mv Login ./bin/gxlogin



