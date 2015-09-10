package GxNet

import (
	"errors"
	"github.com/golang/protobuf/proto"
	"strconv"
	"time"
)

import (
	. "GxMessage"
	. "GxMisc"
	. "GxProto"
	. "GxStatic"
)

type LoginMessageCallback func(*GxTcpConn, *LoginInfo, *GxMessage)

var loginCmds map[uint16]LoginMessageCallback
var gates map[uint32]*GxTcpConn
var loginCounter *Counter

func init() {
	gates = make(map[uint32]*GxTcpConn)
	loginCmds = make(map[uint16]LoginMessageCallback)
	loginCounter = NewCounter()
}

func RegisterLoginMessageCallback(cmd uint16, cb LoginMessageCallback) {
	_, ok := loginCmds[cmd]
	if ok {
		return
	} else {
		loginCmds[cmd] = cb
	}
}

func gateRun(conn *GxTcpConn, gateId uint32) {
	t := time.NewTicker(5 * time.Second)
	go func(conn *GxTcpConn, t *time.Ticker) {
		for {
			select {
			case <-t.C:
				if !conn.Connected {
					return
				}
				SendPbMessage(conn, false, 0, CmdHeartbeat, uint16(loginCounter.Genarate()), 0, nil)
			}
		}
	}(conn, t)

	rdClient := PopRedisClient()
	defer PushRedisClient(rdClient)

	for {
		newMsg, err := conn.Recv()
		if err != nil {
			Debug("Recv error: %s", err)
			return
		}

		if newMsg.GetCmd() != CmdHeartbeat {
			Debug("recv pb msg, info: %s", newMsg.String())
		}

		//获取角色登陆状态
		info := new(LoginInfo)
		playerName := GetGateLoginInfo(rdClient, gateId, newMsg.GetId())
		if playerName != "" {
			info.Get(rdClient, playerName)
		}
		info.GateId = gateId
		info.ConnId = newMsg.GetId()
		info.PlayerName = playerName

		if newMsg.GetCmd() != CmdHeartbeat {
			Debug("init logininfo, gateId: %d, connid: %d", gateId, newMsg.GetId(), *info)
		}

		//没有拉取过角色列表，不能进行接下来的操作
		if newMsg.GetCmd() > CmdGetRoleList && info.PlayerName == "" {
			Debug("player has not select role, gateId: %d, connid: %d, cmd: %d, PlayerName: %s",
				gateId, newMsg.GetId(), newMsg.GetCmd, info.PlayerName)
			SendPbMessage(conn, false, newMsg.GetId(), newMsg.GetCmd(), newMsg.GetSeq(), RetNotLogin, nil)
			continue
		}

		//没有登陆过某个角色列表，不能进行接下来的操作
		if newMsg.GetCmd() > CmdSelectRole && info.RoleId == 0 {
			Debug("player has not select role, gateId: %d, connid: %d, cmd: %d, roleid: %d",
				gateId, newMsg.GetId(), newMsg.GetCmd, info.RoleId)
			SendPbMessage(conn, false, newMsg.GetId(), newMsg.GetCmd(), newMsg.GetSeq(), RetNotLogin, nil)
			continue
		}

		cb, ok := loginCmds[newMsg.GetCmd()]
		if ok {
			go cb(conn, info, newMsg)
		} else {
			SendPbMessage(conn, false, newMsg.GetId(), newMsg.GetCmd(), newMsg.GetSeq(), RetMessageNotSupport, nil)
			Debug("message has not been registered, msg: %s", newMsg.String())
		}
	}
}

func ConnectAllGate() error {
	t := time.NewTicker(10 * time.Second)
	if len(loginCmds) == 0 {
		return errors.New("no cmd is registered")
	}

	var req ServerConnectGateReq

	id, _ := Config.Get("server").Get("id").Int()
	req.Id = proto.Uint32(uint32(id))

	for k, _ := range loginCmds {
		req.Cmds = append(req.Cmds, uint32(k))
	}

	Debug("req: %s", req.String())

	f := func(pb proto.Message) error {
		rdClient := PopRedisClient()
		defer PushRedisClient(rdClient)

		var gatesinfo []*GateInfo
		err := GetAllGate(rdClient, &gatesinfo)
		if err != nil {
			return err
		}

		for i := 0; i < len(gatesinfo); i++ {
			conn, ok := gates[gatesinfo[i].Id]
			if ok {
				if conn.Connected {
					continue
				}
				delete(gates, gatesinfo[i].Id)
			}

			conn = NewTcpConn()
			err = conn.Connect(gatesinfo[i].Host2 + ":" + strconv.Itoa(int(gatesinfo[i].Port2)))
			if err != nil {
				Debug("connnect gate fail, remote: %s:%d", gatesinfo[i].Host2, gatesinfo[i].Port2)
				continue
			}
			Debug("connnect gate ok, remote: %s:%d", gatesinfo[i].Host2, gatesinfo[i].Port2)

			SendPbMessage(conn, false, 0, CmdServerConnectGate, uint16(loginCounter.Genarate()), 0, pb)

			msg, err2 := conn.Recv()
			if err2 != nil || msg.GetRet() != RetSucc {
				Debug("connnect gate fail, remote: %s:%d", gatesinfo[i].Host2, gatesinfo[i].Port2)
				continue
			}

			gates[gatesinfo[i].Id] = conn
			go gateRun(conn, gatesinfo[i].Id)
		}
		return nil
	}

	//先连接一次
	err := f(&req)
	if err != nil {
		return err
	}

	//后面10秒检查一次
	for {
		select {
		case <-t.C:
			err = f(&req)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
