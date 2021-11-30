package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/thesky9531/lareina/log"
	"reflect"
)

func (c *Cache) DelKey(ctx context.Context, key string) {
	conn, err := c.pool.GetContext(ctx)
	if err != nil {
		log.ErrLog("", fmt.Errorf("获取redis conn失败 key(%s),error(%v)", key, err))
		return
	}
	defer conn.Close()
	delKey(conn, key)
}

func (c *Cache) DelMultiKey(ctx context.Context, keys ...string) {
	conn, err := c.pool.GetContext(ctx)
	if err != nil {
		log.ErrLog("", fmt.Errorf("获取redis conn失败 key(%v),error(%v)", keys, err))
		return
	}
	defer conn.Close()
	for _, key := range keys {
		delKey(conn, key)
	}
}

func (c *Cache) RegexpDelKey(ctx context.Context, key string) {
	conn, err := c.pool.GetContext(ctx)
	if err != nil {
		log.ErrLog("", fmt.Errorf("获取redis conn失败 key(%s),error(%v)", key, err))
		return
	}
	defer conn.Close()
	regexpDelKey(conn, key)
}

func (c *Cache) RegexpDelMultiKey(ctx context.Context, keys ...string) {
	conn, err := c.pool.GetContext(ctx)
	if err != nil {
		log.ErrLog("", fmt.Errorf("获取redis conn失败 keys(%v),error(%v)", keys, err))
		return
	}
	defer conn.Close()
	for _, key := range keys {
		regexpDelKey(conn, key)
	}
}

func (c *Cache) GetRawMessage(ctx context.Context, key string) (json.RawMessage, error) {
	conn, err := c.pool.GetContext(ctx)
	if err != nil {
		log.ErrLog("", fmt.Errorf("获取redis conn失败 key(%s),error(%v)", key, err))
		return nil, err
	}
	defer conn.Close()
	return getKeyBytes(conn, key, c.conf.ExpireTime)
}

func (c *Cache) SetRawMessage(ctx context.Context, key string, data json.RawMessage) error {
	conn, err := c.pool.GetContext(ctx)
	if err != nil {
		log.ErrLog("", fmt.Errorf("获取redis conn失败 key(%s),error(%v)", key, err))
		return err
	}
	defer conn.Close()
	return setKeyBytes(conn, key, data, c.conf.ExpireTime)
}

//使用string 整体存储对象
func (c *Cache) GetObject(ctx context.Context, key string, obj interface{}) error {
	conn, err := c.pool.GetContext(ctx)
	if err != nil {
		log.ErrLog("", fmt.Errorf("获取redis conn失败 key(%s),error(%v)", key, err))
		return err
	}
	defer conn.Close()
	data, err := getKeyBytes(conn, key, c.conf.ExpireTime)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &obj)
}

func (c *Cache) SetObject(ctx context.Context, key string, obj interface{}) error {
	conn, err := c.pool.GetContext(ctx)
	if err != nil {
		log.ErrLog("", fmt.Errorf("获取redis conn失败 key(%s),error(%v)", key, err))
		return err
	}
	defer conn.Close()
	data, err := json.Marshal(obj)
	if err != nil {
		log.ErrLog("", fmt.Errorf("mashal obj fail, error(%v)", err))
		return err
	}
	return setKeyBytes(conn, key, data, c.conf.ExpireTime)
}

//使用hash 分字段存储对象
func (c *Cache) GetHashObject(ctx context.Context, key string, fields []string, obj interface{}) error {
	conn, err := c.pool.GetContext(ctx)
	if err != nil {
		log.ErrLog("", fmt.Errorf("获取redis conn失败 key(%s),error(%v)", key, err))
		return err
	}
	defer conn.Close()
	data, err := getKeyHash(conn, key, fields, c.conf.ExpireTime)
	if err != nil {
		return err
	}
	return unmarshalRedisObj(data, reflect.ValueOf(obj))
}

func (c *Cache) SetHashObject(ctx context.Context, key string, fields []string, obj interface{}) error {
	conn, err := c.pool.GetContext(ctx)
	if err != nil {
		log.ErrLog("", fmt.Errorf("获取redis conn失败 key(%s),error(%v)", key, err))
		return err
	}
	defer conn.Close()
	data, err := marshalRedisObj(reflect.ValueOf(obj))
	if err != nil {
		log.ErrLog("", fmt.Errorf("mashal obj fail, error(%v)", err))
		return err
	}
	return setKeyHash(conn, key, fields, data, c.conf.ExpireTime)
}

//使用zset存储id列表
func (c *Cache) GetIdSet(ctx context.Context, key string) ([]int64, error) {
	conn, err := c.pool.GetContext(ctx)
	if err != nil {
		log.ErrLog("", fmt.Errorf("获取redis conn失败 key(%s),error(%v)", key, err))
		return nil, err
	}
	defer conn.Close()
	return getIdSet(conn, key, c.conf.ExpireTime)
}

func (c *Cache) SetIdSet(ctx context.Context, key string, list []int64) error {
	conn, err := c.pool.GetContext(ctx)
	if err != nil {
		log.ErrLog("", fmt.Errorf("获取redis conn失败 key(%s),error(%v)", key, err))
		return err
	}
	defer conn.Close()
	return setIdSet(conn, key, list, c.conf.ExpireTime)
}

//使用Hash存储Name,value List列表
func (c *Cache) GetNameList(ctx context.Context, listKey string) ([]KeyValue, error) {
	conn, err := c.pool.GetContext(ctx)
	if err != nil {
		log.ErrLog("", fmt.Errorf("获取redis conn失败 key(%s),error(%v)", listKey, err))
		return nil, err
	}
	defer conn.Close()
	return getNameList(conn, listKey, c.conf.ExpireTime)
}

func (c *Cache) SetNameList(ctx context.Context, listKey string, list []KeyValue) error {
	conn, err := c.pool.GetContext(ctx)
	if err != nil {
		log.ErrLog("", fmt.Errorf("获取redis conn失败 key(%s),error(%v)", listKey, err))
		return err
	}
	defer conn.Close()
	return setNameList(conn, listKey, list, c.conf.ExpireTime)
}

func (c *Cache) SetKeyInt64List(ctx context.Context, listKey string, list int64) error {
	conn, err := c.pool.GetContext(ctx)
	if err != nil {
		log.ErrLog("", fmt.Errorf("获取redis conn失败 key(%s),error(%v)", listKey, err))
		return err
	}
	defer conn.Close()
	return setKeyInt64(conn, listKey, list, c.conf.ExpireTime)
}

func (c *Cache) GetKeyInt64List(ctx context.Context, key string) (int64, error) {
	conn, err := c.pool.GetContext(ctx)
	if err != nil {
		log.ErrLog("", fmt.Errorf("获取redis conn失败 key(%s),error(%v)", key, err))
		return 0, err
	}
	defer conn.Close()
	return getKeyInt64(conn, key, c.conf.ExpireTime)
}

//在线人数push
func (c *Cache) RPushOnlineCount(ctx context.Context, key string, count int64) (err error) {
	conn, err := c.pool.GetContext(ctx)
	if err != nil {
		log.ErrLog("", fmt.Errorf("获取redis conn失败 key(%s),error(%v)", key, err))
		return err
	}
	defer conn.Close()

	return rpushOnlineCount(conn, key, count)
}

//在线人数数量
func (c *Cache) GetLenOnlineCount(ctx context.Context, key string) (data []int64, err error) {
	conn, err := c.pool.GetContext(ctx)
	if err != nil {
		log.ErrLog("", fmt.Errorf("获取redis conn失败 key(%s),error(%v)", key, err))
		return data, err
	}
	defer conn.Close()
	return getLenOnlineCount(conn, key)
}

//获取当前时间的数量
//func (c *Cache) GetCurrentOnlineCount(ctx context.Context,key string )(data []string, err error){
//	conn, err := c.pool.GetContext(ctx)
//	if err != nil {
//		log.ErrLog("", fmt.Errorf("获取redis conn失败 key(%s),error(%v)", key, err))
//		return data,err
//	}
//	defer conn.Close()
//	return getLenOnlineCount(conn,key)
//}

//list的rpush
func (c *Cache) RPushList(ctx context.Context, key string, id int64) (err error) {
	conn, err := c.pool.GetContext(ctx)
	if err != nil {
		log.ErrLog("", fmt.Errorf("获取redis conn失败 key(%s),error(%v)", key, err))
		return err
	}
	defer conn.Close()

	return rpushList(conn, key, id)
}

//将当前用户的id加入到set中
func (c *Cache) SetSetID(ctx context.Context, key string, id int64) (err error) {
	conn, err := c.pool.GetContext(ctx)
	if err != nil {
		log.ErrLog("", fmt.Errorf("获取redis conn失败 key(%s),error(%v)", key, err))
		return err
	}
	defer conn.Close()

	return setSetID(conn, key, id)
}

func (c *Cache) Ping(ctx context.Context) error {
	conn, err := c.pool.GetContext(ctx)
	if err != nil {
		log.ErrLog("", fmt.Errorf("获取redis conn失败 error(%v)", err))
		return err
	}
	defer conn.Close()
	return nil
}

//获取当前时间的在线数量
func (c *Cache) GetSetCount(ctx context.Context, key string) (count int64, err error) {
	conn, err := c.pool.GetContext(ctx)
	if err != nil {
		log.ErrLog("", fmt.Errorf("获取redis conn失败 key(%s),error(%v)", key, err))
		return 0, err
	}
	defer conn.Close()
	return getSetCount(conn, key)
}

func (c *Cache) GetRegexpKeys(ctx context.Context, key string) ([]string, error) {
	conn, err := c.pool.GetContext(ctx)
	if err != nil {
		log.ErrLog("", fmt.Errorf("获取redis conn失败 key(%s),error(%v)", key, err))
		return nil, err
	}
	defer conn.Close()
	return getRegexpKey(conn, key)
}

func (c *Cache) SetExpireTimeKey(ctx context.Context, key string, value string, expireTime int) error {
	conn, err := c.pool.GetContext(ctx)
	if err != nil {
		log.ErrLog("", fmt.Errorf("获取redis conn失败 key(%s),error(%v)", key, err))
		return err
	}
	defer conn.Close()
	return setKeyString(conn, key, value, expireTime)
}

func (c *Cache) GetExpireTimeKey(ctx context.Context, key string, expireTime int) (string, error) {
	conn, err := c.pool.GetContext(ctx)
	if err != nil {
		log.ErrLog("", fmt.Errorf("获取redis conn失败 key(%s),error(%v)", key, err))
		return "", err
	}
	defer conn.Close()
	return getKeyString(conn, key, expireTime)
}

func (c *Cache) HReset(ctx context.Context, key string, field string) error {
	conn, err := c.pool.GetContext(ctx)
	if err != nil {
		log.ErrLog("", fmt.Errorf("获取redis conn失败 key(%s),error(%v)", key, err))
		return err
	}
	defer conn.Close()
	return hReset(conn, key, field, c.conf.ExpireTime)
}

func (c *Cache) HIncrBy(ctx context.Context, key string, field string, num int64) error {
	conn, err := c.pool.GetContext(ctx)
	if err != nil {
		log.ErrLog("", fmt.Errorf("获取redis conn失败 key(%s),error(%v)", key, err))
		return err
	}
	defer conn.Close()
	return hIncrBy(conn, key, field, num, c.conf.ExpireTime)
}

func (c *Cache) HGetNum(ctx context.Context, key string, field string) (int64, error) {
	conn, err := c.pool.GetContext(ctx)
	if err != nil {
		log.ErrLog("", fmt.Errorf("获取redis conn失败 key(%s),error(%v)", key, err))
		return 0, err
	}
	defer conn.Close()
	return hGetNum(conn, key, field, c.conf.ExpireTime)
}

func (c *Cache) GetInt64(ctx context.Context, key string) (int64, error) {
	conn, err := c.pool.GetContext(ctx)
	if err != nil {
		log.ErrLog("", fmt.Errorf("获取redis conn失败 key(%s),error(%v)", key, err))
		return 0, err
	}
	defer conn.Close()
	return getKeyInt64(conn, key, c.conf.ExpireTime)
}

func (c *Cache) SetInt64(ctx context.Context, key string, data int64) error {
	conn, err := c.pool.GetContext(ctx)
	if err != nil {
		log.ErrLog("", fmt.Errorf("获取redis conn失败 key(%s),error(%v)", key, err))
		return err
	}
	defer conn.Close()
	return setKeyInt64(conn, key, data, c.conf.ExpireTime)
}

func (c *Cache) HDel(ctx context.Context, key string, fields ...string) error {
	conn, err := c.pool.GetContext(ctx)
	if err != nil {
		log.ErrLog("", fmt.Errorf("获取redis conn失败 key(%s),error(%v)", key, err))
		return err
	}
	defer conn.Close()
	return hDel(conn, key, fields)
}
