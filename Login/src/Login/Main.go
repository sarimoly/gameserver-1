package main

import (
	"errors"
	"math/rand"
	"strconv"
	"time"
)

import (
	. "GxMessage"
	. "GxMisc"
	. "GxNet"
	. "GxStatic"
)

var server *GxTcpServer
var r *rand.Rand

func NewConn(conn *GxTcpConn) {
	Debug("new connnect, remote: %s", conn.Conn.RemoteAddr().String())
}

func DisConn(conn *GxTcpConn) {
	Debug("dis connnect, remote: %s", conn.Conn.RemoteAddr().String())
}

func NewMessage(conn *GxTcpConn, msg *GxMessage) error {
	Debug("new message, remote: %s", conn.Conn.RemoteAddr().String())
	conn.Send(msg)
	return errors.New("close")
}

func start_server() {
	port, _ := Config.Get("server").Get("port").Int()
	server = NewGxTcpServer(NewConn, DisConn, NewMessage, true)
	server.RegisterClientCmd(CmdLogin, login)
	server.RegisterClientCmd(CmdRegister, register)
	server.Start(":" + strconv.Itoa(port))
}

func main() {
	LoadConfig("config.json")
	InitLogger("login")

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
	//
	start_server()
	Debug("connect redis fail")
}
