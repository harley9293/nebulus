package httpd

import (
	"github.com/harley9293/nebulus/pkg/db"
	"time"
)

type SessionType int

const (
	SessionTypeLocal = iota
	SessionTypeRedis
)

type Config struct {
	SType       SessionType
	SExpireTime time.Duration
	Redis       *db.RedisConfig
}

func (c *Config) Fill() {
	if c.SExpireTime == 0 {
		c.SExpireTime = 24 * time.Hour
	}
}
