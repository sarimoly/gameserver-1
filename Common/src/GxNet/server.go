package GxNet

import (
	"github.com/golang/protobuf/proto"
	"net"
	"sync"
)

import (
	. "GxMessage"
	. "GxMisc"
	. "GxStatic"
)

type NewConnCallback func(*GxTcpConn)                         //新连接回调
type DisConnCallback func(*GxTcpConn)                         //连接断开回调
type clientMessageCallback func(*GxTcpConn, *GxMessage) error //已经注册的命令回调
type RawMessageCallback func(*GxTcpConn, *GxMessage) error    //没有注册的消息的回调

type GxTcpServer struct {
	mutex     *sync.Mutex
	addrConns map[string]*GxTcpConn //key为连接地址的map
	idConns   map[uint32]*GxTcpConn //key为连接id的map
	id        uint32                //用来给客户端连接分配id

	Nc         NewConnCallback
	Dc         DisConnCallback
	Rm         RawMessageCallback
	clientCmds map[uint16]clientMessageCallback

	MessageCtrl bool
}

func NewGxTcpServer(nc NewConnCallback, dc DisConnCallback, rm RawMessageCallback, messageCtrl bool) *GxTcpServer {
	server := new(GxTcpServer)
	server.addrConns = make(map[string]*GxTcpConn)
	server.idConns = make(map[uint32]*GxTcpConn)
	server.id = 0
	server.mutex = new(sync.Mutex)

	server.Nc = nc
	server.Dc = dc
	server.Rm = rm
	server.clientCmds = make(map[uint16]clientMessageCallback)
	server.MessageCtrl = messageCtrl

	//注册心跳回调
	server.RegisterClientCmd(CmdHeartbeat, HeartbeatCallback)
	return server
}

//注册需要直接处理的消息
func (server *GxTcpServer) RegisterClientCmd(cmd uint16, cb clientMessageCallback) {
	_, ok := server.clientCmds[cmd]
	if ok {
		return
	} else {
		server.clientCmds[cmd] = cb
	}
}

//服务端启动函数
func (server *GxTcpServer) Start(port string) error {

	listener, err := net.Listen("tcp", port)
	if err != nil {
		Debug("lister %s fail", port)
		return err
	}
	Debug("server start, host: %s", port)

	for {
		conn, err1 := listener.Accept()
		if err1 != nil {
			Debug("server Accept fail, err: ", err1)
			return err1
		}
		if conn != nil {
			gxConn := NewTcpConn()
			gxConn.Conn = conn
			gxConn.Connected = true
			gxConn.MessageCtrl = server.MessageCtrl

			server.mutex.Lock()
			//分配id
			server.id++
			gxConn.Id = server.id

			server.idConns[server.id] = gxConn
			server.addrConns[conn.RemoteAddr().String()] = gxConn
			server.mutex.Unlock()

			server.Nc(gxConn)

			go server.runConn(gxConn)
		}
	}

	return nil
}

func (server *GxTcpServer) runConn(gxConn *GxTcpConn) {
	go gxConn.runHeartbeat()

	for {
		var msg *GxMessage = nil
		var err error

		empty := gxConn.GetProcessMsgNull()
		if empty {
			//没有消息正在处理，查找当前缓存的消息
			msg = gxConn.GetUnprocessMsg()
		}
		if msg == nil {
			//如果没有缓存的消息，则读取新信息
			msg, err = gxConn.Recv()
			if err != nil {
				server.closeConn(gxConn)
				return
			}
			if !empty {
				//如果有消息正在处理，缓存刚刚收到的消息
				Debug("client has processing message, remote: %s, msg: %s", gxConn.Conn.RemoteAddr().String(),
					msg.String())
				gxConn.SaveUnprocessMsg(msg)
				continue
			}
		}

		if msg.GetCmd() != CmdHeartbeat {
			Debug("recv buff msg, info: %s", msg.String())
		}

		if cb, ok := server.clientCmds[msg.GetCmd()]; ok {
			//消息已经被注册
			err = cb(gxConn, msg)
			if err != nil {
				//回调返回值不为空，则关闭连接
				server.closeConn(gxConn)
				return
			}
			continue
		}

		//消息没有被注册
		err = server.Rm(gxConn, msg)
		if err != nil {
			server.closeConn(gxConn)
			return
		}
	}
}

func (server *GxTcpServer) closeConn(gxConn *GxTcpConn) {
	server.Dc(gxConn)
	server.mutex.Lock()
	delete(server.addrConns, gxConn.Conn.RemoteAddr().String())
	delete(server.idConns, gxConn.Id)
	server.mutex.Unlock()
	gxConn.Toc <- 0xFFFF
	gxConn.Conn.Close()
}

func (server *GxTcpServer) FindConnByRetome(retome string) *GxTcpConn {
	server.mutex.Lock()
	defer server.mutex.Unlock()
	info, ok := server.addrConns[retome]
	if ok {
		return info
	} else {
		return nil
	}
}

func (server *GxTcpServer) FindConnById(id uint32) *GxTcpConn {
	server.mutex.Lock()
	defer server.mutex.Unlock()
	info, ok := server.idConns[id]
	if ok {
		return info
	} else {
		return nil
	}
}

func HeartbeatCallback(conn *GxTcpConn, msg *GxMessage) error {
	conn.Toc <- 0
	msg.SetRet(RetSucc)
	conn.Send(msg)
	return nil
}

func SendRawMessage(conn *GxTcpConn, notify bool, id uint32, cmd uint16, seq uint16, ret uint16, buff []byte) {
	msg := NewGxMessage()
	msg.SetId(id)
	msg.SetCmd(cmd)
	msg.SetRet(ret)
	msg.SetSeq(seq)
	msg.SetMask(MessageMaskNotify, notify)

	if len(buff) == 0 {
		msg.SetLen(0)
		msg.SetUnlen(0)
	} else {
		err := msg.Package(buff)
		if err != nil {
			Debug("PackagePbmsg error")
			return
		}
	}

	conn.Send(msg)
	if msg.GetCmd() != CmdHeartbeat {
		Debug("send buff msg, info: %s", string(buff))
	}
}

func SendPbMessage(conn *GxTcpConn, notify bool, id uint32, cmd uint16, seq uint16, ret uint16, pb proto.Message) {
	msg := NewGxMessage()
	msg.SetId(id)
	msg.SetCmd(cmd)
	msg.SetRet(ret)
	msg.SetSeq(seq)
	msg.SetMask(MessageMaskNotify, notify)

	if pb == nil {
		msg.SetLen(0)
		msg.SetUnlen(0)
	} else {
		err := msg.PackagePbmsg(pb)
		if err != nil {
			Debug("PackagePbmsg error")
			return
		}
	}

	if msg.GetCmd() != CmdHeartbeat {
		Debug("send buff msg, info: %s", msg.String())
	}
	conn.Send(msg)
}
