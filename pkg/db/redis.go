package db

import (
	"github.com/go-redis/redis"
	log "github.com/harley9293/blotlog"
)

type RedisConfig struct {
	Addr string
	Pass string
	DB   int
}

type RedisClient struct {
	*redis.Client
}

func NewRedisClient(config *RedisConfig) *RedisClient {
	client := &RedisClient{
		Client: redis.NewClient(&redis.Options{
			Addr:     config.Addr,
			Password: config.Pass,
			DB:       config.DB,
		}),
	}

	_, err := client.Ping().Result()
	if err != nil {
		log.Error("redis connect error: %s", err.Error())
		return nil
	}
	return client
}
