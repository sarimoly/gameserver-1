package GxStatic

import (
	"gopkg.in/redis.v3"
	"testing"
)

func Test_newPalayerID(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       7,  // use default DB
	})

	client.Del(PlayerIdTableName)

	id := newPalayerID(client)
	if id != 100001 {
		t.Error("new player id error: ", id)
	}

	id = newPalayerID(client)
	if id != 100002 {
		t.Error("new player id error: ", id)
	}

	id = newPalayerID(client)
	if id != 100003 {
		t.Error("new player id error: ", id)
	}
}

func Test_player(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       7,  // use default DB
	})

	client.Del(PlayerIdTableName)

	player := NewPlayer(client, "username", "password")

	if !player.VerifyPassword("password") {
		t.Error("Verify Password 1 error: ", player)
	}

	if player.VerifyPassword("password123") {
		t.Error("Verify Password 2 error: ", player)
	}

	player.Save(client)

	///////////////////////////////////
	player1 := new(Player)
	err := player1.Get(client, "username")
	if err != nil {
		t.Error("get player 1 error: ", err)
	}
	if player.Id != player1.Id || player.Username != player1.Username || player.Password != player1.Password {
		t.Error("get player 1 error: ", player, player1)
	}

	token := player1.SaveToken(client)
	name := CheckToken(client, token)
	if name != player1.Username {
		t.Error("Check Token 1 error: ", name)
	}
	name = CheckToken(client, "aaaaaaa")
	if name != "" {
		t.Error("Check Token 1 error: ", name)
	}
}
