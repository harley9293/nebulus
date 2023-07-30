package db

import (
	"context"
	"github.com/go-redis/redis"
	log "github.com/harley9293/blotlog"
)

type RedisClient struct {
	*redis.Client
	ctx context.Context
}

func NewRedisClient(Addr, Pass string, DB int) *RedisClient {
	client := &RedisClient{
		Client: redis.NewClient(&redis.Options{
			Addr:     Addr,
			Password: Pass,
			DB:       DB,
		}),
		ctx: context.Background(),
	}

	_, err := client.Ping().Result()
	if err != nil {
		log.Error("redis connect error: %s", err.Error())
		return nil
	}
	return client
}
