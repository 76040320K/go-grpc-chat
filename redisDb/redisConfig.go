package redisDb

import (
	"github.com/go-redis/redis/v8"
)

var (
	RdsCli *redis.Client
)

func RedisClient() {
	RdsCli = redis.NewClient(&redis.Options{
		Addr: "10.186.69.228:6379",
		Password: "",
		DB: 0,
	})
}