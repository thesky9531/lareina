package redis

import (
	"context"
	"github.com/gomodule/redigo/redis"
	"time"
)

type Config struct {
	Network     string
	Addr        string
	MaxIdle     int
	MaxActive   int
	IdleTimeout int
	Wait        bool
	ExpireTime  int
	PassWord    string
}

type KeyValue struct {
	Key   string
	Value string
}

type Cache struct {
	pool *redis.Pool
	conf *Config
}

func New(c *Config) *Cache {
	return &Cache{
		pool: &redis.Pool{
			DialContext: func(ctx context.Context) (conn redis.Conn, e error) {
				conn, e = redis.DialContext(
					ctx,
					c.Network,
					c.Addr,
					redis.DialPassword(c.PassWord),
				)
				return
			},
			MaxIdle:     c.MaxIdle,
			MaxActive:   c.MaxActive,
			IdleTimeout: time.Second * time.Duration(c.IdleTimeout),
			Wait:        c.Wait,
		},
		conf: c,
	}
}

func (c *Cache) Close() {
	c.pool.Close()
}
