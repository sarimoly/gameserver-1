package GxStatic

import (
	"gopkg.in/redis.v3"
	"testing"
)

func Test_GameServer(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       7,  // use default DB
	})

	server1 := &GameServer{
		Id:     1,
		Name:   "server1",
		Status: 1,
	}

	server2 := &GameServer{
		Id:     2,
		Name:   "server2",
		Status: 1,
	}

	err := SaveGameServer(client, server1)
	if err != nil {
		t.Error("SaveGameServer 1 error: ", err)
	}
	err = SaveGameServer(client, server2)
	if err != nil {
		t.Error("SaveGameServer 2 error: ", err)
	}

	var servers []*GameServer
	err = GetAllGameServer(client, &servers)
	if err != nil {
		t.Error("GetAllGameServer  1 error: ", err)
	}

	if len(servers) != 2 {
		t.Error("GetAllGameServer 2 error: ", len(servers))
	}

	if servers[0].Id != 1 && servers[0].Id != 2 {
		t.Error("GetAllGameServer 3 error: ", servers[0])
	}
}
