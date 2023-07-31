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
	options *redis.Options
}

func NewRedisClient(config *RedisConfig) *RedisClient {
	client := &RedisClient{
		options: &redis.Options{
			Addr:     config.Addr,
			Password: config.Pass,
			DB:       config.DB,
		},
	}

	client.connect()
	return client
}

func (c *RedisClient) connect() {
	c.Client = redis.NewClient(c.options)
	_, err := c.Ping().Result()
	if err != nil {
		log.Error("redis connect error: %s", err.Error())
	}
}

func (c *RedisClient) KeepAlive() {
	_, err := c.Ping().Result()
	if err != nil {
		log.Error("redis keepalive error: %s, Try Reconnect", err.Error())
		c.Shutdown()
		c.connect()
	}
}

func (c *RedisClient) Shutdown() {
	err := c.Close()
	if err != nil {
		log.Error("redis shutdown error: %s", err.Error())
	}
}
