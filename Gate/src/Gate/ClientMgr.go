package main

import (
	"errors"
)

import (
	. "GxMessage"
	. "GxMisc"
	. "GxNet"
)

var clientRouter *GxTcpServer

func init() {
	clientRouter = NewGxTcpServer(ClientNewConn, ClientDisConn, ClientRawMessage, true)
}

func ClientNewConn(conn *GxTcpConn) {
	Debug("new client connnect, remote: %s", conn.Conn.RemoteAddr().String())
	addClient()
}

func ClientDisConn(conn *GxTcpConn) {
	Debug("dis client connnect, remote: %s", conn.Conn.RemoteAddr().String())
	subClient()
}

func ClientRawMessage(conn *GxTcpConn, msg *GxMessage) error {
	Debug("new client message, remote: %s %s", conn.Conn.RemoteAddr().String(), msg.String())

	_, server := GetServerByCmd(msg.GetCmd())
	if server == nil {
		Debug("msg is not register, remote: %s, msg: %s", conn.Conn.RemoteAddr().String(),
			msg.String())
		return errors.New("close")
	}
	msg.SetId(conn.Id)
	server.Send(msg)
	return nil
}
