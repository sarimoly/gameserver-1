package GxStatic

import (
	"gopkg.in/redis.v3"
	"strconv"
)

import (
	. "GxMisc"
)

var GateInfoTableName string = "h_gate_info"

type GateInfo struct {
	Id    uint32
	Host1 string
	Port1 uint32
	Host2 string
	Port2 uint32
	Count uint32
	Ts    int64
}

func SaveGate(client *redis.Client, gate *GateInfo) error {
	buf, err := MsgToBuf(gate)
	if err != nil {
		return err
	}

	client.HSet(GateInfoTableName, strconv.Itoa(int(gate.Id)), string(buf))

	return nil
}

func GetAllGate(client *redis.Client, gates *[]*GateInfo) error {
	m := client.HGetAllMap(GateInfoTableName)
	r, err := m.Result()
	if err != nil {
		return err
	}

	for _, v := range r {
		j, err2 := BufToMsg([]byte(v))
		if err2 != nil {
			return err2
		}
		gate := new(GateInfo)
		JsonToStruct(j, gate)
		*gates = append(*gates, gate)
	}
	return nil
}
