package main

import (
	"fmt"
	"gopkg.in/redis.v3"
	"time"
)
import (
	. "GxMisc"
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

type TestStruct struct {
	Uid        int    `PK`
	Username   string `PK`
	Departname string
	Created    int64
}

func main() {
	LoadConfig("config.json")
	InitLogger("NewServer")

	if !connect_redis() {
		Debug("connect redis fail")
		return
	}

	var saveone TestStruct
	saveone.Uid = 1
	saveone.Username = "name"
	saveone.Departname = "Test Add Departname"
	saveone.Created = time.Now().Unix()

	SaveToRedis(rdClient, &saveone)

	var saveone1 TestStruct
	saveone1.Uid = 1
	saveone1.Username = "name"
	LoadFromRedis(rdClient, &saveone1)
	fmt.Println(saveone1)
	fmt.Println("ok")
}
