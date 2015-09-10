package GxStatic

import (
	"gopkg.in/redis.v3"
	"testing"
)

func Test_RoleName(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       7,  // use default DB
	})

	client.Del(RoleNameListTableName)
	client.Del(RoleIdTableName)

	role1 := &Role{
		Id:           NewRoleID(client),
		PlayerName:   "guang1",
		GameServerId: 1,
		Name:         "role1",
		VocationId:   1,
		Level:        1,
		Expr:         1000,
		GodValue:     1000,
		Prestige:     1000,
		Gold:         1000,
		Crystal:      1000,
	}

	role2 := &Role{
		Id:           NewRoleID(client),
		PlayerName:   "guang2",
		GameServerId: 1,
		Name:         "role2",
		VocationId:   2,
		Level:        2,
		Expr:         2000,
		GodValue:     2000,
		Prestige:     2000,
		Gold:         2000,
		Crystal:      2000,
	}

	if role1.Id != 10000001 {
		t.Error("role id 1 error: ", role1.Id)
	}
	if role2.Id != 10000002 {
		t.Error("role id 1 error: ", role2.Id)
	}

	SaveRoleName(client, role1.Name)
	SaveRoleName(client, role2.Name)

	if !CheckRoleNameConflict(client, role1.Name) {
		t.Error("CheckRoleNameConflict 1 error: ", role1.Name)
	}
	if !CheckRoleNameConflict(client, role2.Name) {
		t.Error("CheckRoleNameConflict 2 error: ", role2.Name)
	}
	if CheckRoleNameConflict(client, "role3") {
		t.Error("CheckRoleNameConflict 3 error: ")
	}

	role1.Save(client)

	role3 := &Role{}
	err := role3.Get(client, 10000001)
	if err != nil {
		t.Error("role3 get 1 error: ", err)
	}
	if role3.Id != role1.Id || role3.PlayerName != role1.PlayerName || role3.GameServerId != role1.GameServerId ||
		role3.Name != role1.Name || role3.Level != role1.Level || role3.Expr != role1.Expr {
		t.Error("role3 get 2 error: ", role3)
	}
	err = role3.Get(client, 10000003)
	if err == nil {
		t.Error("role3 get 3 error: ")
	}
}
