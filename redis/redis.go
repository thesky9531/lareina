package redis

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"log"
	"time"
)

var session *RedisSession

type RedisConfig struct {
	Host, Password, Prefix             string
	Port, DbName, MaxIdle, IdleTimeout int
}
type RedisSession struct {
	pool   *redis.Pool
	prefix string
}

func LoadRedisSession(c *RedisConfig) error {
	session = &RedisSession{}

	session.pool = &redis.Pool{
		MaxIdle:     c.MaxIdle,
		IdleTimeout: time.Duration(c.IdleTimeout) * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", fmt.Sprintf("%s:%d", c.Host, c.Port),
				redis.DialPassword(c.Password), redis.DialDatabase(c.DbName))
		},
	}
	if err := session.pool.Get().Err(); err != nil {
		return err
	}
	return nil
}

func GetSession() *RedisSession {
	if session == nil {
		session = &RedisSession{}
	}
	return session
}

func (r *RedisSession) Setprefix(name string) {
	r.prefix = name + ":"
}

func (r *RedisSession) Getprefix() string {
	return r.prefix
}

func (r *RedisSession) GetConn() redis.Conn {
	return r.pool.Get()
}

func Close(conn redis.Conn) {
	func(conn redis.Conn) {
		err := conn.Close()
		if err != nil {
			log.Println("redis conn close failure")
		}
	}(conn)
}

func (r *RedisSession) Do(cmd string, args ...interface{}) (interface{}, error) {
	conn := r.pool.Get()
	defer Close(conn)
	args[0] = r.Getprefix() + args[0].(string)
	return conn.Do(cmd, args...)
}
