package main

import (
	"fmt"
	"gopkg.in/redis.v3"
	"os"
	"strconv"
)
import (
	. "GxMisc"
	. "GxStatic"
)

var rdClient *redis.Client

func connect_redis() bool {
	rdHost, _ := Config.Get("redis").Get("host").String()
	rdPort, _ := Config.Get("redis").Get("port").Int()
	rdDb, _ := Config.Get("redis").Get("db").Int64()
	rdClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", rdHost, rdPort),
		Password: "",   // no password set
		DB:       rdDb, // use default DB
	})
	return rdClient != nil
}

func main() {
	argNum := len(os.Args)
	if argNum != 4 {
		fmt.Println("NewServer <id> <name> <status>")
		fmt.Println("0 - hot")
		fmt.Println("1 - new")
		fmt.Println("2 - maintain")
		return
	}

	LoadConfig("config.json")
	InitLogger("NewServer")

	if !connect_redis() {
		Debug("connect redis fail")
		return
	}

	id, _ := strconv.Atoi(os.Args[1])
	status, _ := strconv.Atoi(os.Args[3])

	server := &GameServer{
		Id:     uint32(id),
		Name:   os.Args[2],
		Status: uint32(status),
	}

	err := SaveGameServer(rdClient, server)
	if err != nil {
		Debug("SaveGameServer 1 error: %s", err)
		return
	}

	fmt.Println("ok")
}
