package GxNet

import (
	"container/list"
	"errors"
	"fmt"
	"net"
	"sync"
	"time"
)

import (
	. "GxMessage"
	. "GxMisc"
)

type GxTcpConn struct {
	Id        uint32 //连接id
	Conn      net.Conn
	Connected bool

	TimeoutCount int          //超时次数
	T            *time.Ticker //超时检测定时器
	Toc          chan int

	Data []interface{} //预留使用

	mutex       *sync.Mutex
	msg         *GxMessage //当前正在处理的消息
	msgList     *list.List //未处理的消息列表
	MessageCtrl bool
}

func NewTcpConn() *GxTcpConn {
	tcpConn := new(GxTcpConn)
	tcpConn.Connected = false
	tcpConn.TimeoutCount = 0
	tcpConn.Toc = make(chan int, 1)
	tcpConn.T = time.NewTicker(5 * time.Second)

	tcpConn.mutex = new(sync.Mutex)
	tcpConn.msg = nil
	tcpConn.msgList = list.New()
	tcpConn.MessageCtrl = false

	return tcpConn
}

func (conn *GxTcpConn) GetUnprocessMsg() *GxMessage {
	if !conn.MessageCtrl {
		return nil
	}

	conn.mutex.Lock()
	defer conn.mutex.Unlock()

	if conn.msgList.Len() == 0 {
		return nil
	}
	msg := conn.msgList.Front().Value.(*GxMessage)
	conn.msgList.Remove(conn.msgList.Front())
	return msg
}

func (conn *GxTcpConn) SaveUnprocessMsg(msg *GxMessage) {
	if !conn.MessageCtrl {
		return
	}

	conn.mutex.Lock()
	defer conn.mutex.Unlock()

	conn.msgList.PushBack(msg)
}

func (conn *GxTcpConn) GetProcessMsgNull() bool {
	if !conn.MessageCtrl {
		return true
	}

	conn.mutex.Lock()
	defer conn.mutex.Unlock()

	return conn.msg == nil
}

func (conn *GxTcpConn) SaveProcessMsg(msg *GxMessage) {
	if !conn.MessageCtrl {
		return
	}

	conn.mutex.Lock()
	defer conn.mutex.Unlock()

	//服务端发送的通知seq为0
	if msg.GetMask(MessageMaskNotify) {
		return
	}

	if conn.msg == nil {
		conn.msg = msg
	} else {
		if conn.msg.GetSeq() == msg.GetSeq() && conn.msg.GetCmd() == msg.GetCmd() {
			conn.msg = nil
		}
	}

}

//处理心跳函数，用协程启动
func (conn *GxTcpConn) runHeartbeat() {
	for {
		select {
		case state := <-conn.Toc:
			if state == 0XFFFF {
				return
			}
			conn.TimeoutCount = state
		case <-conn.T.C:
			if conn.TimeoutCount > 3 {
				//超时超过三次关闭连接
				conn.Conn.Close()

				Debug("client[%d] %s timeout", conn.Id, conn.Conn.RemoteAddr().String())
				return
			} else if conn.TimeoutCount >= 0 {
				conn.TimeoutCount = conn.TimeoutCount + 1
			} else {
				break
			}
		}
	}
}

//发送消息
func (conn *GxTcpConn) Send(msg *GxMessage) error {
	//发送结束设置当前正在处理的消息为nil，通知消息忽略
	conn.SaveProcessMsg(msg)

	//读取消息头
	len, err := conn.Conn.Write(msg.Header)
	if err != nil {
		fmt.Println(err)
	}
	if uint16(len) != MessageHeaderLen {
		return errors.New("send error")
	}

	//如果消息体没有数据，直接返回
	if msg.GetLen() == 0 {
		return nil
	}

	//读取消息体
	len, err = conn.Conn.Write(msg.Data)
	if err != nil {
		fmt.Println(err)
	}
	if uint16(len) != msg.GetLen() {
		return errors.New("send error")
	}
	return nil
}

func (conn *GxTcpConn) Recv() (*GxMessage, error) {
	//写消息头
	msg := NewGxMessage()
	len, err := conn.Conn.Read(msg.Header)
	if err != nil {
		conn.Connected = false
		return nil, err
	}
	if uint16(len) != MessageHeaderLen {
		return nil, errors.New("recv error")
	}

	//消息头没有数据，则返回
	if msg.GetLen() == 0 {
		return msg, nil
	}

	//写消息体
	msg.InitData()
	len, err = conn.Conn.Read(msg.Data)
	if err != nil {
		conn.Connected = false
		return nil, err
	}
	if uint16(len) != msg.GetLen() {
		return nil, errors.New("recv error")
	}
	return msg, nil
}

//连接指定host
func (conn *GxTcpConn) Connect(host string) error {
	c, err := net.Dial("tcp", host)

	if err != nil {
		return err
	}
	conn.Conn = c
	conn.Connected = true
	return nil
}
