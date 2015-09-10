package GxStatic

import (
	"gopkg.in/redis.v3"
	"testing"
	"time"
)

func Test_login_info(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       7,  // use default DB
	})

	info := new(LoginInfo)
	info.GateId = 1
	info.ConnId = 123
	info.BeginTs = time.Now().Unix()
	info.EndTs = time.Now().Unix() + 10
	info.ServerId = 1
	info.RoleId = 1000001
	info.PlayerName = "guang"
	info.Save(client)

	info2 := new(LoginInfo)
	err := info2.Get(client, "guang")

	if err != nil {
		t.Error("login info get 1 error: ", err)
	}
	if info.BeginTs != info2.BeginTs {
		t.Error("login info get 2 error: ", info, info2)
	}

	err = info.Get(client, "guang1")
	if err == nil {
		t.Error("login info get 3 error: ", err)
	}

	//
	SaveGateLoginInfo(client, info.GateId, info.ConnId, "guang")
	name := GetGateLoginInfo(client, info.GateId, info.ConnId)
	if name != "guang" {
		t.Error("GetGateLoginInfo error: ", name)
	}
}
