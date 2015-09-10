package GxStatic

import (
	"crypto/md5"
	// "errors"
	"fmt"
	"gopkg.in/redis.v3"
	"io"
	// "strconv"
	"sync"
	"time"
)

import (
	. "GxMisc"
)

var salt1 = "f4g*h(j"
var salt2 = "1^2&4*(d)"

var PlayerTableName = "h_player:"
var PlayerIdTableName = "k_player_id"
var PlayerTokenTableName = "k_player_token:"

var idMutex *sync.Mutex

type Player struct {
	Id       uint32
	Username string `PK`
	Password string
}

func init() {
	idMutex = new(sync.Mutex)
}

func newPalayerID(client *redis.Client) uint32 {
	idMutex.Lock()
	defer idMutex.Unlock()

	if !client.Exists(PlayerIdTableName).Val() {
		client.Set(PlayerIdTableName, "100000", 0)

	}
	return uint32(client.Incr(PlayerIdTableName).Val())
}

func generatePassward(username string, password string) string {
	h := md5.New()

	io.WriteString(h, salt1)
	io.WriteString(h, username)
	io.WriteString(h, salt2)
	io.WriteString(h, password)

	return fmt.Sprintf("%x", h.Sum(nil))
}

func NewPlayer(client *redis.Client, username string, password string) *Player {
	player := &Player{
		Id:       newPalayerID(client),
		Username: username,
		Password: generatePassward(username, password),
	}
	return player
}

func (player *Player) checkUsernameConflict(client *redis.Client) bool {
	return client.Exists(PlayerTableName + player.Username).Val()
}

func (player *Player) Save(client *redis.Client) {
	SaveToRedis(client, player)
	// client.HMSet(PlayerTableName+player.Username, "id", strconv.Itoa(int(player.Id)), "username", player.Username, "password", player.Password)
}

func (player *Player) Get(client *redis.Client, username string) error {
	// key := PlayerTableName + username

	// if !client.Exists(key).Val() {
	// 	return errors.New("role is not exists")
	// }

	// id, err := client.HGet(key, "id").Uint64()
	// if err != nil {
	// 	return err
	// }
	// player.Id = uint32(id)
	// player.Username = client.HGet(key, "username").Val()
	// player.Password = client.HGet(key, "password").Val()
	player.Username = username
	return LoadFromRedis(client, player)
}

func (player *Player) VerifyPassword(password string) bool {
	return player.Password == generatePassward(player.Username, password)
}

func (player *Player) SaveToken(client *redis.Client) string {
	h := md5.New()
	io.WriteString(h, player.Username)
	io.WriteString(h, fmt.Sprintf("%ld", time.Now().Unix()))
	io.WriteString(h, "2@#RR#R@")
	token := fmt.Sprintf("%x", h.Sum(nil))

	client.Set(PlayerTokenTableName+token, player.Username, time.Hour*1)
	return token
}

func CheckToken(client *redis.Client, token string) string {
	key := PlayerTokenTableName + token
	if !client.Exists(key).Val() {
		return ""
	}
	return client.Get(key).Val()
}
