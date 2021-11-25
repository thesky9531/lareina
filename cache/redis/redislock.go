package redis

import (
	"context"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/satori/go.uuid"
	"github.com/thesky9531/lareina/log"
	"strconv"
	"time"
)

var delScript = redis.NewScript(1, `
if redis.call("get", KEYS[1]) == ARGV[1] then
	return redis.call("del", KEYS[1])
else
	return 0
end`)

const lockEx int = 10

type Lock struct {
	key   string
	token uuid.UUID
}

func getRedisKey(key string) string {
	return fmt.Sprintf("redislock_%s", key)
}

/*
   redis 类型 字符串设置一个分布式锁 (哈希内部字段不支持过期判断,redis只支持顶级key过期)
   @param key: 锁名,格式为  用户id_操作_方法
   @param requestId:  客户端唯一id 用来指定锁不被其他线程(协程)删除
   @param ex: 过期时间
*/
func (c *Cache) Lock(ctx context.Context, key string) *Lock {
	select {
	case <-ctx.Done():
		return nil
	default:
	}
	lock := c.addLock(ctx, key)
	for lock == nil {
		time.Sleep(100 * time.Millisecond)
		select {
		case <-ctx.Done():
			return nil
		default:
		}
		lock = c.addLock(ctx, key)
	}
	return lock
}

func (c *Cache) addLock(ctx context.Context, key string) *Lock {
	conn, err := c.pool.GetContext(ctx)
	if err != nil {
		log.ErrLog("", fmt.Errorf("获取redis conn失败 key(%s),error(%v)", key, err))
		return nil
	}
	defer conn.Close()
	token := uuid.NewV4()
	redisKey := getRedisKey(key)
	msg, _ := redis.String(
		conn.Do("SET", redisKey, token, SetIfNotExist, SetWithExpireTime, lockEx),
	)
	if msg == SetLockSuccess {
		return &Lock{
			key:   redisKey,
			token: token,
		}
	}
	return nil
}

/*
   删除redis分布式锁

   @param key:redis类型字符串的key值
   @param requestId: 唯一值id,与value值对比,避免在分布式下其他实例删除该锁
*/
func (c *Cache) Unlock(ctx context.Context, l *Lock) {
	if l == nil {
		return
	}
	conn, err := c.pool.GetContext(ctx)
	if err != nil {
		log.ErrLog("", fmt.Errorf("获取redis conn失败 error(%v)", err))
		return
	}
	defer conn.Close()
	msg, err := redis.Int64(delScript.Do(conn, l.key, l.token))
	// 避免操作时间过长,自动过期时再删除返回结果为0
	if err != nil {
		log.ErrLog(strconv.FormatInt(msg, 10), err)
	}
}
