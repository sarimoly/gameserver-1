package main

import (
	"math/rand"
	"time"
)

import (
	. "GxMisc"
	. "GxNet"
)

var r *rand.Rand

func main() {
	LoadConfig("config.json")
	InitLogger("center")

	r = rand.New(rand.NewSource(time.Now().UnixNano()))

	//
	rdHost, _ := Config.Get("redis").Get("host").String()
	rdPort, _ := Config.Get("redis").Get("port").Int()
	rdDb, _ := Config.Get("redis").Get("db").Int64()
	err := ConnectRedis(rdHost, rdPort, rdDb)
	if err != nil {
		Debug("connect redis fail, err: %s", err)
		return
	}

	err = ConnectAllGate()
	if err != nil {
		Debug("ConnectAllGate fail, %s", err)
		return
	}
}
