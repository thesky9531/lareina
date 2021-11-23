package redis

import (
	"fmt"
	"testing"
)

func TestRids(t *testing.T) {
	rc := RedisConfig{
		Host:        "139.196.232.xxx",
		Password:    "xxxredis...xxx",
		Prefix:      "test",
		Port:        6379,
		DbName:      15,
		MaxIdle:     20,
		IdleTimeout: 240,
	}
	err := LoadRedisSession(rc)
	if err != nil {
		fmt.Println(err)
	}
	se := GetSession()
	re, err := se.Do("set", "test", "test111")
	fmt.Println(re, err)
	re1, err := se.Do("get", "test")
	fmt.Println(re1, err)
}
