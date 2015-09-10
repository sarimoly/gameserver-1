package main

import (
	"errors"
	"github.com/golang/protobuf/proto"
	"time"
)

import (
	. "GxMessage"
	. "GxMisc"
	. "GxNet"
	. "GxProto"
	. "GxStatic"
)

func login(conn *GxTcpConn, msg *GxMessage) error {
	rdClient := PopRedisClient()
	defer PushRedisClient(rdClient)

	var req LoginServerReq
	var rsp LoginServerRsp
	err := msg.UnpackagePbmsg(&req)
	if err != nil {
		Debug("UnpackagePbmsg error")
		return errors.New("close")
	}

	player := new(Player)
	err = player.Get(rdClient, req.GetRaw().GetUsername())
	if err != nil {
		Debug("user is not exists, username: %s", req.GetRaw().GetUsername())
		SendPbMessage(conn, false, 0, CmdLogin, msg.GetSeq(), RetUserNotExists, nil)
		return errors.New("close")
	}

	if !player.VerifyPassword(req.GetRaw().GetPwd()) {
		SendPbMessage(conn, false, 0, CmdLogin, msg.GetSeq(), RetPwdError, nil)
		return errors.New("close")
	}

	Debug("old user: %s login from %s", req.GetRaw().GetUsername(), conn.Conn.RemoteAddr().String())

	rsp.Info = &LoginRspInfo{
		Token: proto.String(player.SaveToken(rdClient)),
	}

	var gates []*GateInfo
	GetAllGate(rdClient, &gates)
	var gate *GateInfo = nil
	for i := 0; i < len(gates); i++ {
		if (time.Now().Unix() - gates[i].Ts) > 10 {
			continue
		}
		if gate == nil {
			gate = gates[i]
		} else {
			if gate.Count > gates[i].Count {
				gate = gates[i]
			}
		}
	}
	if gate != nil {
		rsp.GetInfo().Host = proto.String(gate.Host1)
		rsp.GetInfo().Port = proto.Uint32(gate.Port1)
	}
	// 	optional uint32 index = 1; //区号.
	// optional string name = 2; //服务器名称.
	// optional uint32 statue = 3; //服务器状态.
	// optional uint32 lastts = 4; //最近登录时间.

	var servers []*GameServer
	GetAllGameServer(rdClient, &servers)
	for i := 0; i < len(servers); i++ {
		rsp.GetInfo().Srvs = append(rsp.GetInfo().Srvs, &GameSrvInfo{
			Index:  proto.Uint32(servers[i].Id),
			Name:   proto.String(servers[i].Name),
			Status: proto.Uint32(servers[i].Status),
			Lastts: proto.Uint32(0),
		})
	}

	SendPbMessage(conn, false, 0, CmdLogin, msg.GetSeq(), RetSucc, &rsp)
	return errors.New("close")
}

func register(conn *GxTcpConn, msg *GxMessage) error {
	rdClient := PopRedisClient()
	defer PushRedisClient(rdClient)

	var req LoginServerReq
	err := msg.UnpackagePbmsg(&req)
	if err != nil {
		Debug("UnpackagePbmsg error")
		return errors.New("close")
	}

	player := new(Player)
	err = player.Get(rdClient, req.Raw.GetUsername())
	if err == nil {
		Debug("user has been exists, username: %s", req.GetRaw().GetUsername())
		SendPbMessage(conn, false, 0, CmdRegister, msg.GetSeq(), RetUserExists, nil)
		return errors.New("close")
	}

	player = NewPlayer(rdClient, req.Raw.GetUsername(), req.GetRaw().GetPwd())
	player.Save(rdClient)

	Debug("new user: %s login from %s", req.Raw.GetUsername(), conn.Conn.RemoteAddr().String())
	SendPbMessage(conn, false, 0, CmdRegister, msg.GetSeq(), RetSucc, &LoginServerRsp{
		Info: &LoginRspInfo{
			Token: proto.String(player.SaveToken(rdClient)),
		},
	})
	return errors.New("close")
}
