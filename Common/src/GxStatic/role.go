package GxStatic

import (
	// "errors"
	// "fmt"
	"gopkg.in/redis.v3"
	"strconv"
	"sync"
)

import (
	. "GxMisc"
)

var RoleTableName = "h_role:"
var RoleListTableName = "l_role_list:"

var RoleNameListTableName = "s_role_name"
var RoleIdTableName = "k_role_id"

var roleIdMutex *sync.Mutex

type Role struct {
	Id           uint32 `PK`
	PlayerName   string
	GameServerId uint32
	Name         string
	VocationId   uint32 //职业
	Level        uint32
	Expr         uint64
	GodValue     uint64 //神格
	Prestige     uint64 //声望
	Gold         uint64 //金币
	Crystal      uint64 //水晶
}

func init() {
	roleIdMutex = new(sync.Mutex)
}

func NewRoleID(client *redis.Client) uint32 {
	roleIdMutex.Lock()
	defer roleIdMutex.Unlock()

	if !client.Exists(RoleIdTableName).Val() {
		client.Set(RoleIdTableName, "10000000", 0)

	}
	return uint32(client.Incr(RoleIdTableName).Val())
}

func (role *Role) Save(client *redis.Client) {
	gameServerId := strconv.Itoa(int(role.GameServerId))
	id := strconv.Itoa(int(role.Id))
	client.LPush(RoleListTableName+role.PlayerName+":"+gameServerId, id)

	SaveToRedis(client, role)

	// client.HMSet(RoleTableName+id, "id", id, "player_name", role.PlayerName, "game_server_id", strconv.Itoa(int(role.GameServerId)),
	// 	"name", role.Name, "vocation_id", strconv.Itoa(int(role.VocationId)),
	// 	"level", fmt.Sprintf("%d", role.Level), "expr", fmt.Sprintf("%d", role.Expr),
	// 	"god_value", fmt.Sprintf("%d", role.GodValue), "prestige", fmt.Sprintf("%d", role.Prestige),
	// 	"gold", fmt.Sprintf("%d", role.Gold), "crystal", fmt.Sprintf("%d", role.Crystal))
}

func (role *Role) Get(client *redis.Client, id uint32) error {
	role.Id = id
	return LoadFromRedis(client, role)

	// key := RoleTableName + strconv.Itoa(int(id))

	// if !client.Exists(key).Val() {
	// 	return errors.New("role is not exists")
	// }

	// i, _ := strconv.Atoi(client.HGet(key, "id").Val())
	// role.Id = uint32(i)
	// role.PlayerName = client.HGet(key, "player_name").Val()
	// i, _ = strconv.Atoi(client.HGet(key, "game_server_id").Val())
	// role.GameServerId = uint32(i)
	// role.Name = client.HGet(key, "name").Val()
	// i, _ = strconv.Atoi(client.HGet(key, "vocation_id").Val())
	// role.VocationId = uint32(i)
	// i, _ = strconv.Atoi(client.HGet(key, "level").Val())
	// role.Level = uint32(i)
	// role.Expr, _ = client.HGet(key, "expr").Uint64()

	// role.GodValue, _ = client.HGet(key, "god_value").Uint64()
	// role.Prestige, _ = client.HGet(key, "prestige").Uint64()
	// role.Gold, _ = client.HGet(key, "gold").Uint64()
	// role.Crystal, _ = client.HGet(key, "crystal").Uint64()
	// return nil
}

func GetRoleList(client *redis.Client, playerName string, gameServerId uint32) []string {
	return client.LRange(RoleListTableName+playerName+":"+strconv.Itoa(int(gameServerId)), 0, -1).Val()
}

func CheckRoleNameConflict(client *redis.Client, name string) bool {
	return client.SIsMember(RoleNameListTableName, name).Val()
}

func SaveRoleName(client *redis.Client, name string) {
	client.SAdd(RoleNameListTableName, name)
}
