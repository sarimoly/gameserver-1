package GxStatic

import (
	"gopkg.in/redis.v3"
	"testing"
)

func Test_gate(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       7,  // use default DB
	})

	client.Del(GateInfoTableName)

	gate := &GateInfo{
		Id:    1,
		Host1: "127.0.0.1",
		Port1: 10000,
		Host2: "127.0.0.1",
		Port2: 10001,
		Count: 100,
		Ts:    100,
	}
	err := SaveGate(client, gate)
	if err != nil {
		t.Error("gate save error: ", err)
	}

	var gates []*GateInfo
	err = GetAllGate(client, &gates)
	if err != nil {
		t.Error("get all gate  1 error: ", err)
	}

	if len(gates) != 1 {
		t.Error("gate all gate 2 error: ", len(gates))
	}

	if gates[0].Id != 1 || gates[0].Host1 != "127.0.0.1" || gates[0].Port1 != 10000 {
		t.Error("gate all gate 3 error: ", gates[0])
	}

	gate2 := &GateInfo{
		Id:    2,
		Host1: "127.0.0.1",
		Port1: 10002,
		Host2: "127.0.0.1",
		Port2: 10003,
		Count: 100,
		Ts:    100,
	}
	err = SaveGate(client, gate2)
	if err != nil {
		t.Error("gate save error: ", err)
	}

	var gates2 []*GateInfo
	err = GetAllGate(client, &gates2)
	if err != nil {
		t.Error("get all gate  1 error: ", err)
	}

	if len(gates2) != 2 {
		t.Error("gate all gate 2 error: ", len(gates2))
	}

}
