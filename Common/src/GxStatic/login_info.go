package GxStatic

import (
	// "errors"
	// "fmt"
	"gopkg.in/redis.v3"
	"strconv"
)

import (
	. "GxMisc"
)

type LoginInfo struct {
	GateId     uint32
	ConnId     uint32
	BeginTs    int64
	EndTs      int64
	ServerId   uint32
	RoleId     uint32
	PlayerName string `PK`
}

var LoginInfoTableName = "h_login_info:"
var GateLoginInfoTableName = "k_gate_login_info:"

func (info *LoginInfo) Save(client *redis.Client) {
	SaveToRedis(client, info)
	// key := LoginInfoTableName + info.PlayerName

	// client.HMSet(key, "gate_id", strconv.Itoa(int(info.GateId)), "conn_id", strconv.Itoa(int(info.ConnId)),
	// 	"begin_ts", fmt.Sprintf("%d", info.BeginTs), "end_ts", fmt.Sprintf("%d", info.EndTs),
	// 	"server_id", strconv.Itoa(int(info.ServerId)), "role_id", strconv.Itoa(int(info.RoleId)),
	// 	"player_name", info.PlayerName)
}

func (info *LoginInfo) Get(client *redis.Client, palyerName string) error {
	// key := LoginInfoTableName + palyerName

	// if !client.Exists(key).Val() {
	// 	return errors.New("LoginInfo is not exists")
	// }

	// i, _ := client.HGet(key, "gate_id").Int64()
	// info.GateId = uint32(i)
	// i, _ = client.HGet(key, "conn_id").Int64()
	// info.ConnId = uint32(i)

	// info.BeginTs, _ = client.HGet(key, "begin_ts").Int64()
	// info.EndTs, _ = client.HGet(key, "end_ts").Int64()

	// i, _ = client.HGet(key, "server_id").Int64()
	// info.ServerId = uint32(i)
	// i, _ = client.HGet(key, "role_id").Int64()
	// info.RoleId = uint32(i)

	// info.PlayerName = client.HGet(key, "player_name").Val()
	info.PlayerName = palyerName
	return LoadFromRedis(client, info)
}

func SaveGateLoginInfo(client *redis.Client, gateId uint32, connId uint32, palyerName string) {
	key := GateLoginInfoTableName + strconv.Itoa(int(gateId)) + ":" + strconv.Itoa(int(connId))

	client.Set(key, palyerName, 0)
}

func GetGateLoginInfo(client *redis.Client, gateId uint32, connId uint32) string {
	key := GateLoginInfoTableName + strconv.Itoa(int(gateId)) + ":" + strconv.Itoa(int(connId))

	return client.Get(key).Val()
}

func DelGateLoginInfo(client *redis.Client, gateId uint32, connId uint32) {
	key := GateLoginInfoTableName + strconv.Itoa(int(gateId)) + ":" + strconv.Itoa(int(connId))

	client.Del(key)
}

func DisconnLogin(client *redis.Client, token string, info *LoginInfo) uint16 {
	playerName := CheckToken(client, token)
	if playerName == "" {
		return RetTokenError
	}
	oldInfo := new(LoginInfo)
	oldInfo.Get(client, playerName)

	if info.GateId != oldInfo.GateId || info.ConnId != oldInfo.ConnId {
		info.PlayerName = oldInfo.PlayerName
		info.ServerId = oldInfo.ServerId
		info.RoleId = oldInfo.RoleId
		info.BeginTs = oldInfo.BeginTs
		DelGateLoginInfo(client, oldInfo.GateId, oldInfo.ConnId)
		SaveGateLoginInfo(client, info.GateId, info.ConnId, playerName)
		info.Save(client)
	}
	return RetSucc
}
