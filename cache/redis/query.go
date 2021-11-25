package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/thesky9531/lareina/log"
	"reflect"
)

func (c *Cache) QueryRawMessage(ctx context.Context, key string,
	update func() (json.RawMessage, error)) (rsp json.RawMessage, err error) {
	//首先判断缓存中有没有
	//首先从缓存获取
	rsp, err = c.GetRawMessage(ctx, key)
	if err == nil {
		return rsp, err
	}
	//缓存没有或失败，从db获取
	lock := c.Lock(ctx, key)
	defer c.Unlock(ctx, lock)
	//再次判断是否有数据
	rsp, err = c.GetRawMessage(ctx, key)
	if err == nil {
		return rsp, err
	}
	rsp, err = update()
	if err != nil {
		return rsp, err
	}
	err = c.SetRawMessage(ctx, key, rsp)
	if err != nil {
	}
	return rsp, err
}

func (c *Cache) QueryHashObject(ctx context.Context, key string, fields []string, obj interface{},
	update func() ([]string, interface{}, error)) error {
	conn, err := c.pool.GetContext(ctx)
	if err != nil {
		log.ErrLog("", fmt.Errorf("获取redis conn失败 key(%s),error(%v)", key, err))
		return err
	}
	defer conn.Close()
	//首先判断缓存中有没有
	//首先从缓存获取
	data, err := getKeyHash(conn, key, fields, c.conf.ExpireTime)
	if err == nil {
		err = unmarshalRedisObj(data, reflect.ValueOf(obj))
		if err == nil {
			return err
		}
	}
	//缓存没有或失败，从db获取
	lock := c.Lock(ctx, key)
	defer c.Unlock(ctx, lock)
	//再次判断是否有数据
	data, err = getKeyHash(conn, key, fields, c.conf.ExpireTime)
	if err == nil {
		err = unmarshalRedisObj(data, reflect.ValueOf(obj))
		if err == nil {
			return err
		}
	}
	//获取更新数据
	var uptData interface{}
	fields, uptData, err = update()
	if err != nil {
		return err
	}
	data, err = marshalRedisObj(reflect.ValueOf(uptData))
	if err != nil {
		return err
	}
	err = setKeyHash(conn, key, fields, data, c.conf.ExpireTime)
	if err != nil {
	}
	return unmarshalRedisObj(data, reflect.ValueOf(obj))
}

func (c *Cache) QueryIdSet(ctx context.Context, key string,
	update func() ([]int64, error)) (rsp []int64, err error) {
	//首先判断缓存中有没有
	//首先从缓存获取
	rsp, err = c.GetIdSet(ctx, key)
	if err == nil {
		return rsp, err
	}
	//缓存没有或失败，从db获取
	lock := c.Lock(ctx, key)
	defer c.Unlock(ctx, lock)
	//再次判断是否有数据
	rsp, err = c.GetIdSet(ctx, key)
	if err == nil {
		return rsp, err
	}
	rsp, err = update()
	if err != nil {
		return rsp, err
	}
	err = c.SetIdSet(ctx, key, rsp)
	if err != nil {
	}
	return rsp, err
}

func (c *Cache) QueryNameList(ctx context.Context, listKey string,
	update func() ([]KeyValue, error)) (rsp []KeyValue, err error) {
	//首先判断缓存中有没有
	//首先从缓存获取
	rsp, err = c.GetNameList(ctx, listKey)
	if err == nil {
		return rsp, err
	}
	//缓存没有或失败，从db获取
	lock := c.Lock(ctx, listKey)
	defer c.Unlock(ctx, lock)
	//再次判断是否有数据
	rsp, err = c.GetNameList(ctx, listKey)
	if err == nil {
		return rsp, err
	}
	rsp, err = update()
	if err != nil {
		return rsp, err
	}
	err = c.SetNameList(ctx, listKey, rsp)
	return rsp, err
}
