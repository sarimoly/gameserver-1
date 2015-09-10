package main

import (
	"github.com/golang/protobuf/proto"
	"strconv"
	"time"
)

import (
	. "GxMessage"
	. "GxMisc"
	. "GxNet"
	. "GxProto"
	. "GxStatic"
)

func init() {
	RegisterLoginMessageCallback(CmdHeartbeat, ClientHeartbeatCallback)
	RegisterLoginMessageCallback(CmdGetRoleList, GetRoleListCallback)
	RegisterLoginMessageCallback(CmdSelectRole, SelectRoleCallback)
	RegisterLoginMessageCallback(CmdCreateRole, CreateRoleCallback)
}

func ClientHeartbeatCallback(conn *GxTcpConn, info *LoginInfo, msg *GxMessage) {

}

func GetRoleListCallback(conn *GxTcpConn, info *LoginInfo, msg *GxMessage) {
	rdClient := PopRedisClient()
	defer PushRedisClient(rdClient)

	var req GetRoleListReq
	var rsp GetRoleListRsp
	err := msg.UnpackagePbmsg(&req)
	if err != nil {
		SendPbMessage(conn, false, msg.GetId(), msg.GetCmd(), msg.GetSeq(), RetFail, nil)
		return
	}
	if req.Info == nil || req.GetInfo().Token == nil || req.ServerId == nil {
		SendPbMessage(conn, false, msg.GetId(), msg.GetCmd(), msg.GetSeq(), RetMsgFormatError, nil)
		return
	}

	playerName := CheckToken(rdClient, req.GetInfo().GetToken())
	if playerName == "" {
		SendPbMessage(conn, false, msg.GetId(), msg.GetCmd(), msg.GetSeq(), RetTokenError, nil)
		return
	}

	player := new(Player)
	err = player.Get(rdClient, playerName)
	if err != nil {
		SendPbMessage(conn, false, msg.GetId(), msg.GetCmd(), msg.GetSeq(), RetUserNotExists, nil)
		return
	}

	ids := GetRoleList(rdClient, playerName, req.GetServerId())
	for i := 0; i < len(ids); i++ {
		id, _ := strconv.Atoi(ids[i])

		role := new(Role)
		err = role.Get(rdClient, uint32(id))
		if err != nil {
			Debug("role %d is not existst", id)
			continue
		}
		rsp.Roles = append(rsp.Roles, &RoleCommonInfo{
			Id:         proto.Uint32(role.Id),
			Name:       proto.String(role.Name),
			Level:      proto.Uint32(role.Level),
			VocationId: proto.Uint32(role.VocationId),
		})
	}

	info.PlayerName = playerName
	info.BeginTs = time.Now().Unix()
	info.ServerId = req.GetServerId()
	SaveGateLoginInfo(rdClient, info.GateId, info.ConnId, playerName)
	info.Save(rdClient)

	SendPbMessage(conn, false, msg.GetId(), msg.GetCmd(), msg.GetSeq(), RetSucc, &rsp)
}

func SelectRoleCallback(conn *GxTcpConn, info *LoginInfo, msg *GxMessage) {
	rdClient := PopRedisClient()
	defer PushRedisClient(rdClient)

	var req SelectRoleReq

	err := msg.UnpackagePbmsg(&req)
	if err != nil {
		SendPbMessage(conn, false, msg.GetId(), msg.GetCmd(), msg.GetSeq(), RetFail, nil)
		return
	}

	if req.RoleId == nil {
		SendPbMessage(conn, false, msg.GetId(), msg.GetCmd(), msg.GetSeq(), RetMsgFormatError, nil)
		return
	}

	if req.Info != nil && req.GetInfo().Token != nil {
		//重新重连
		ret := DisconnLogin(rdClient, req.GetInfo().GetToken(), info)
		if ret != RetSucc {
			SendPbMessage(conn, false, msg.GetId(), msg.GetCmd(), msg.GetSeq(), ret, nil)
			return
		}
	}

	role := new(Role)
	err = role.Get(rdClient, req.GetRoleId())
	if err != nil {
		Debug("role %d is not existst", req.GetRoleId())
		SendPbMessage(conn, false, msg.GetId(), msg.GetCmd(), msg.GetSeq(), RetRoleNotExists, nil)
		return
	}

	info.RoleId = req.GetRoleId()
	info.Save(rdClient)

	SendPbMessage(conn, false, msg.GetId(), msg.GetCmd(), msg.GetSeq(), RetSucc, &SelectRoleRsp{
		Role: &RoleCommonInfo{
			Id:         proto.Uint32(role.Id),
			Name:       proto.String(role.Name),
			Level:      proto.Uint32(role.Level),
			VocationId: proto.Uint32(role.VocationId),
			Expr:       proto.Uint64(role.Expr),
			GodValue:   proto.Uint64(role.GodValue),
			Prestige:   proto.Uint64(role.Prestige),
			Gold:       proto.Uint64(role.Gold),
			Crystal:    proto.Uint64(role.Crystal),
		},
	})

}

func CreateRoleCallback(conn *GxTcpConn, info *LoginInfo, msg *GxMessage) {
	rdClient := PopRedisClient()
	defer PushRedisClient(rdClient)

	var req CreateRoleReq

	err := msg.UnpackagePbmsg(&req)
	if err != nil {
		SendPbMessage(conn, false, msg.GetId(), msg.GetCmd(), msg.GetSeq(), RetFail, nil)
		return
	}

	if req.Name == nil || req.VocationId == nil {
		SendPbMessage(conn, false, msg.GetId(), msg.GetCmd(), msg.GetSeq(), RetMsgFormatError, nil)
		return
	}

	if req.Info != nil && req.GetInfo().Token != nil {
		//重新重连
		ret := DisconnLogin(rdClient, req.GetInfo().GetToken(), info)
		if ret != RetSucc {
			SendPbMessage(conn, false, msg.GetId(), msg.GetCmd(), msg.GetSeq(), ret, nil)
			return
		}
	}

	ids := GetRoleList(rdClient, info.PlayerName, info.ServerId)
	if len(ids) > 0 {
		SendPbMessage(conn, false, msg.GetId(), msg.GetCmd(), msg.GetSeq(), RetRoleExists, nil)
		return
	}

	if CheckRoleNameConflict(rdClient, req.GetName()) {
		SendPbMessage(conn, false, msg.GetId(), msg.GetCmd(), msg.GetSeq(), RetRoleNameConflict, nil)
		return
	}

	role := &Role{
		Id:           NewRoleID(rdClient),
		PlayerName:   info.PlayerName,
		GameServerId: info.ServerId,
		Name:         req.GetName(),
		VocationId:   req.GetVocationId(),
		Level:        0,
		Expr:         0,
		GodValue:     0,
		Prestige:     0,
		Gold:         10000,
		Crystal:      0,
	}
	role.Save(rdClient)
	SaveRoleName(rdClient, role.Name)
	//
	info.RoleId = role.Id
	info.Save(rdClient)
	//
	SendPbMessage(conn, false, msg.GetId(), msg.GetCmd(), msg.GetSeq(), RetSucc, &CreateRoleRsp{
		Role: &RoleCommonInfo{
			Id:         proto.Uint32(role.Id),
			Name:       proto.String(role.Name),
			Level:      proto.Uint32(role.Level),
			VocationId: proto.Uint32(role.VocationId),
			Expr:       proto.Uint64(role.Expr),
			GodValue:   proto.Uint64(role.GodValue),
			Prestige:   proto.Uint64(role.Prestige),
			Gold:       proto.Uint64(role.Gold),
			Crystal:    proto.Uint64(role.Crystal),
		},
	})
}
