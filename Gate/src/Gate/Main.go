package main

import (
	"math/rand"
	"strconv"
	"time"
)

import (
	. "GxMisc"
)

var r *rand.Rand

func start_server() {
	host2, _ := Config.Get("server").Get("host2").String()
	port1, _ := Config.Get("server").Get("port1").Int()
	port2, _ := Config.Get("server").Get("port2").Int()

	go clientRouter.Start(host2 + ":" + strconv.Itoa(port1))

	serverRouter.Start(host2 + ":" + strconv.Itoa(port2))
}

func main() {
	LoadConfig("config.json")
	InitLogger("gate")

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

	gate_run()
	//
	start_server()
}
