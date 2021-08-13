package redis

import (
	"github.com/go-redis/redis/v8"
	"github.com/midaef/fibonacci-service/config"
)

func NewRedisClient(config *config.Config) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.Redis.Host + ":" + config.Redis.Port,
		Password: config.Redis.Password,
		DB:       config.Redis.DB,
	})

	return rdb
}
