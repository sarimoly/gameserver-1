package GxMisc

import (
	"container/list"
	"errors"
	"fmt"
	"gopkg.in/redis.v3"
	"sync"
)

var reidsClients *list.List
var reidsMutex *sync.Mutex

var redisHost string
var redisPort int
var redisDb int64

var redisCount int

func init() {
	reidsClients = list.New()
	reidsMutex = new(sync.Mutex)
	redisCount = 4
}

func ConnectRedis(host string, port int, db int64) error {
	redisHost = host
	redisPort = port
	redisDb = db

	for i := 0; i < redisCount; i++ {
		rdClient := redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", redisHost, redisPort),
			Password: "",      // no password set
			DB:       redisDb, // use default DB
		})
		if rdClient == nil {
			return errors.New("connect redis fail")
		}
		reidsClients.PushBack(rdClient)
	}
	return nil
}

func PopRedisClient() *redis.Client {
	reidsMutex.Lock()
	defer reidsMutex.Unlock()
	if reidsClients.Len() == 0 {
		for i := 0; i < redisCount; i++ {
			rdClient := redis.NewClient(&redis.Options{
				Addr:     fmt.Sprintf("%s:%d", redisHost, redisPort),
				Password: "",      // no password set
				DB:       redisDb, // use default DB
			})
			if rdClient == nil {
				return nil
			}
			reidsClients.PushBack(rdClient)
		}
		redisCount += redisCount
	}

	client := reidsClients.Front().Value.(*redis.Client)
	reidsClients.Remove(reidsClients.Front())
	return client
}

func PushRedisClient(client *redis.Client) {
	reidsMutex.Lock()
	defer reidsMutex.Unlock()

	reidsClients.PushBack(client)
}
