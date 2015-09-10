package main

import (
	"sync"
)

import (
	. "GxMessage"
	. "GxMisc"
	. "GxNet"
	. "GxProto"
	. "GxStatic"
)

type ServersInfo map[uint32]*GxTcpConn

var serverRouter *GxTcpServer
var mutex *sync.Mutex
var CmdServers map[uint16]ServersInfo //保存所有注册cmd的服务器

func init() {
	mutex = new(sync.Mutex)
	CmdServers = make(map[uint16]ServersInfo)

	serverRouter = NewGxTcpServer(ServerNewConn, ServerDisConn, ServerRawMessage, false)
	serverRouter.RegisterClientCmd(CmdServerConnectGate, ServerConnectGateCallback)
}

func ServerNewConn(conn *GxTcpConn) {
	Debug("new server connnect, remote: %s", conn.Conn.RemoteAddr().String())
}

func ServerDisConn(conn *GxTcpConn) {
	Debug("dis server connnect, remote: %s", conn.Conn.RemoteAddr().String())

	mutex.Lock()
	for i := 0; i < len(conn.Data); i++ {
		delete(CmdServers[uint16(conn.Data[i].(uint32))], conn.Id)
	}
	mutex.Unlock()
}

func ServerRawMessage(conn *GxTcpConn, msg *GxMessage) error {
	Debug("new server message, remote: %s %s", conn.Conn.RemoteAddr().String(), msg.String())

	client := clientRouter.FindConnById(msg.GetId())
	if client == nil {
		Debug("msg cannot find target, remote: %s, msg: %s", conn.Conn.RemoteAddr().String(),
			msg.String())
		return nil
	}

	client.Send(msg)
	if msg.GetMask(MessageMaskDisconn) {
		client.Conn.Close()
	}
	return nil
}

func ServerConnectGateCallback(conn *GxTcpConn, msg *GxMessage) error {
	//register server
	var req ServerConnectGateReq
	err := msg.UnpackagePbmsg(&req)
	if err != nil {
		SendPbMessage(conn, false, 0, msg.GetCmd(), msg.GetSeq(), RetFail, nil)
		return err
	}

	conn.Data = make([]interface{}, len(req.Cmds))
	mutex.Lock()
	for i := 0; i < len(req.Cmds); i++ {
		//保存自己处理的消息cmd
		conn.Data[i] = req.Cmds[i]

		cmd := uint16(req.Cmds[i])
		cmdinfo, ok := CmdServers[cmd]
		if !ok {
			CmdServers[cmd] = make(ServersInfo)
			cmdinfo, _ = CmdServers[cmd]
		}
		cmdinfo[req.GetId()] = conn
	}
	mutex.Unlock()

	SendPbMessage(conn, false, 0, msg.GetCmd(), msg.GetSeq(), RetSucc, nil)
	return nil
}

func GetServerByCmd(cmd uint16) (uint32, *GxTcpConn) {
	info, ok := CmdServers[cmd]
	if ok {
		count := len(info)
		var v []uint32
		for i, _ := range info {
			v = append(v, i)
		}
		id := v[r.Intn(count)]
		conn, ok1 := info[id]
		if ok1 {
			return id, conn
		}
	}
	return 0, nil
}
