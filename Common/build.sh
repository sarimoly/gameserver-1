#!/bin/bash

CURDIR=`pwd`
export GOPATH=$GOPATH':'$CURDIR

echo $GOPATH
go test ./src/GxMessage
go build ./src/GxMessage
go install ./src/GxMessage

go test ./src/GxMisc
go build ./src/GxMisc
go install ./src/GxMisc

go test ./src/GxStatic
go build ./src/GxNet
go install ./src/GxNet

go test ./src/GxStatic
go build ./src/GxStatic
go install ./src/GxStatic

#build proto file
cd ./src/GxProto
./make_proto.sh
cd ../../

go test ./src/GxProto
go build ./src/GxProto
go install ./src/GxProto

echo '==========>build ok'



