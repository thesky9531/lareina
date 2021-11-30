package redis

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/gomodule/redigo/redis"
	"github.com/thesky9531/lareina/log"
)

const (
	SetIfNotExist     = "NX" // 不存在则执行
	SetWithExpireTime = "EX" // 过期时间(秒)  PX 毫秒
	SetLockSuccess    = "OK" // 操作成功
)

//删除key
func delKey(conn redis.Conn, key string) {
	if _, err := conn.Do("DEL", key); err != nil {
		err = fmt.Errorf("delKey:%v;key=%s", err, key)
		log.ErrLog("", err)
	}
}

//模糊删除key
func regexpDelKey(conn redis.Conn, key string) {
	keys := make([]string, 0)
	var err error
	if keys, err = redis.Strings(conn.Do("Keys", key)); err != nil {
		err = fmt.Errorf("regexpDelKey:%v;key=%s", err, key)
		log.ErrLog("", err)
		return
	}
	if len(keys) == 0 {
		return
	}
	args := make([]interface{}, 0)
	for _, k := range keys {
		args = append(args, k)
	}
	if _, err = conn.Do("DEL", args...); err != nil {
		err = fmt.Errorf("redisDelKey:%v;key=%s", err, key)
		log.ErrLog("", err)
	}
}

func getKeyBytes(conn redis.Conn, key string,
	expireTime int) ([]byte, error) {
	var data []byte
	ok, err := redis.Bool(conn.Do("EXPIRE", key, expireTime))
	if err != nil && err != redis.ErrNil {
		err = fmt.Errorf("getKeyBytes conn.Do(EXPIRE, %s) error(%v)", key, err)
		log.ErrLog("", err)
		return data, err
	}
	if !ok {
		err = redis.ErrNil
		return data, err
	}

	data, err = redis.Bytes(conn.Do("GET", key))
	if err != nil && err != redis.ErrNil {
		err = fmt.Errorf("getKeyBytes conn.Do(GET, %s) error(%v)", key, err)
		log.ErrLog("", err)
		return data, err
	}
	return data, nil
}

func setKeyBytes(conn redis.Conn, key string,
	data []byte, expireTime int) error {
	_, err := conn.Do("SET", key, data, SetWithExpireTime, expireTime)
	if err != nil {
		err = fmt.Errorf("RedisSetKeyBytes conn.Do(SET, %v) error(%v)", key, err)
		log.ErrLog("", err)
	}
	return err
}

func getKeyInt64(conn redis.Conn, key string,
	expireTime int) (int64, error) {
	var data int64
	ok, err := redis.Bool(conn.Do("EXPIRE", key, expireTime))
	if err != nil && err != redis.ErrNil {
		err = fmt.Errorf("getKeyInt64 conn.Do(EXPIRE, %s) error(%v)", key, err)
		log.ErrLog("", err)
		return data, err
	}
	if !ok {
		err = redis.ErrNil
		return data, err
	}

	data, err = redis.Int64(conn.Do("GET", key))
	if err != nil && err != redis.ErrNil {
		err = fmt.Errorf("getKeyBytes conn.Do(GET, %s) error(%v)", key, err)
		log.ErrLog("", err)
		return data, err
	}
	return data, nil
}

func setKeyInt64(conn redis.Conn, key string,
	data int64, expireTime int) error {
	_, err := conn.Do("SET", key, data, SetWithExpireTime, expireTime)
	if err != nil {
		err = fmt.Errorf("RedisSetKeyBytes conn.Do(SET, %v) error(%v)", key, err)
		log.ErrLog("", err)
	}
	return err
}

func getKeyHash(conn redis.Conn, key string,
	fields []string, expireTime int) (map[string][]byte, error) {
	data := make(map[string][]byte)
	ok, err := redis.Bool(conn.Do("EXPIRE", key, expireTime))
	if err != nil && err != redis.ErrNil {
		err = fmt.Errorf("getKeyInfo conn.Do(EXPIRE, %s) error(%v)", key, err)
		log.ErrLog("", err)
		return data, err
	}
	if !ok {
		err = redis.ErrNil
		return data, err
	}
	args := make([]interface{}, 0)
	args = append(args, key)
	for _, v := range fields {
		args = append(args, v)
	}
	reply, err := redis.ByteSlices(conn.Do("HMGET", args...))
	if err != nil && err != redis.ErrNil {
		if e, ok := err.(*redis.Error); ok && strings.Index(e.Error(), "WRONGTYPE") == -1 {
			err = fmt.Errorf("redisGetKeyInfo conn.Do(HMGET, %s,%v) error(%v)", key, fields, err)
			log.ErrLog("", err)
			return data, err
		}
		err = nil
		return data, err
	}
	if len(reply) != len(fields) {
		err = errors.New("返回值与查询属性不一致")
		log.ErrLog("", err)
		return data, err
	}
	for i, _ := range fields {
		data[fields[i]] = reply[i]
	}
	return data, nil
}

func setKeyHash(conn redis.Conn, key string, fields []string,
	data map[string][]byte, expireTime int) error {
	if len(fields) <= 0 || len(data) <= 0 {
		// 设置哨兵
		if _, err := conn.Do("SETEX", key, expireTime, "emptylist"); err != nil {
			err = fmt.Errorf("setKeyInfo conn.Do(SETEX, %s,) error(%v)", key, err)
			log.ErrLog("", err)
			return err
		}
		return nil
	}

	args := make([]interface{}, len(fields)*2+1)
	args[0] = key

	for i, f := range fields {
		v, ok := data[f]
		if !ok {
			err := fmt.Errorf("setKeyInfo unsupported field(%v), key(%s)", f, key)
			log.ErrLog("", err)
			return err
		}
		args[i*2+1] = f
		args[i*2+2] = v
	}

	if _, err := conn.Do("HMSET", args...); err != nil {
		err = fmt.Errorf("setKeyInfo conn.Do(HMSET, %v) error(%v)", args, err)
		log.ErrLog("", err)
		return err
	}
	_, err := conn.Do("EXPIRE", key, expireTime)
	if err != nil {
		log.ErrLog("", err)
	}
	return err
}

func marshalRedisObj(dataValue reflect.Value) (rsp map[string][]byte, err error) {
	var v []byte
	rsp = make(map[string][]byte)
	for dataValue.Kind() == reflect.Ptr {
		dataValue = dataValue.Elem()
	}
	if dataValue.Kind() != reflect.Struct {
		err = errors.New("json marshal fail : data not struct or struct pointer")
		log.ErrLog("", err)
		return
	}
	fileNum := dataValue.NumField()
	for j := 0; j < fileNum; j++ {
		f := dataValue.Type().Field(j).Tag.Get("redis")
		if f == "" {
			err = errors.New("json marshal fail : cant get redis tag")
			log.ErrLog("", err)
			return
		}
		v, err = json.Marshal(dataValue.Field(j).Interface())
		if err != nil {
			log.ErrLog("", err)
			return
		}
		rsp[f] = v
	}
	return
}

func unmarshalRedisObj(data map[string][]byte, baseValue reflect.Value) (err error) {
	if baseValue.Kind() != reflect.Ptr {
		err = fmt.Errorf("unmarshal fail %s", baseValue.Kind().String())
		log.ErrLog("", err)
		return
	}
	for baseValue.Kind() == reflect.Ptr {
		baseValue = baseValue.Elem()
	}
	if baseValue.Kind() != reflect.Struct {
		err = fmt.Errorf("unmarshal fail %s", baseValue.Kind().String())
		log.ErrLog("", err)
		return
	}
	fieldNum := baseValue.NumField()
	for f, v := range data {
		if len(v) == 0 || strings.Contains(string(v), "null") {
			continue
		}
		bfind := false
		for i := 0; i < fieldNum; i++ {
			if baseValue.Type().Field(i).Tag.Get("redis") == f {
				err = json.Unmarshal(v, baseValue.Field(i).Addr().Interface())
				if err != nil {
					log.ErrLog("", err)
					return
				}
				bfind = true
				break
			}
		}
		if !bfind {
			err = errors.New(fmt.Sprintf("cant find field %s", f))
			log.ErrLog("", err)
			return

		}
	}
	return nil
}

func getIdSet(conn redis.Conn, listKey string, expireTime int) ([]int64, error) {
	list := make([]int64, 0)
	var listIds []string
	var ListId int64
	listIds, err := getStringSet(conn, listKey, expireTime)
	if err != nil {
		return list, err
	}
	for _, id := range listIds {
		ListId, err = strconv.ParseInt(id, 10, 64)
		if err != nil {
			err = fmt.Errorf("getKeyList ParseInt(%d) error(%v)", ListId, err)
			log.ErrLog("", err)
			return nil, err
		}
		list = append(list, ListId)
	}
	return list, nil
}

//设置redis列表
func setIdSet(conn redis.Conn, listKey string,
	list []int64, expireTime int) error {
	set := make([]string, 0)
	for _, id := range list {
		set = append(set, strconv.FormatInt(id, 10))
	}
	err := setStringSet(conn, listKey, set, expireTime)
	return err
}

func getStringSet(conn redis.Conn, listKey string, expireTime int) ([]string, error) {
	list := make([]string, 0)
	ok, err := redis.Bool(conn.Do("EXPIRE", listKey, expireTime))
	if err != nil && err != redis.ErrNil {
		err = fmt.Errorf("getStringList conn.Do(EXPIRE, %s) error(%v)", listKey, err)
		log.ErrLog("", err)
		return list, err
	}
	if !ok {
		// 不存在或者为空都会重新去DB获取
		err = redis.ErrNil
		return list, err
	}
	list, err = redis.Strings(conn.Do("ZRANGE", listKey, "0", "-1"))
	if err != nil && err != redis.ErrNil {
		if e, ok := err.(*redis.Error); ok && strings.Index(e.Error(), "WRONGTYPE") == -1 {
			err = fmt.Errorf("getStringList conn.Do(SMEMBERS, %s) error(%v)", listKey, err)
			log.ErrLog("", err)
			return list, err
		}
	}
	return list, nil
}

//设置redis列表
func setStringSet(conn redis.Conn,
	listKey string, list []string, expireTime int) (err error) {
	if len(list) <= 0 {
		// 设置哨兵
		if _, err = conn.Do("SETEX", listKey, expireTime, "emptylist"); err != nil {
			err = fmt.Errorf("setStringList conn.Do(SETEX, %s, %v) error(%v)", listKey, list, err)
			log.ErrLog("", err)
		}
		return err
	}
	if _, err = conn.Do("DEL", listKey); err != nil && err != redis.ErrNil {
		err = fmt.Errorf("setStringList conn.Do(DEL, %s) error(%v)", listKey, err)
		log.ErrLog("", err)
		return err
	}
	args := make([]interface{}, 2*len(list)+1)
	args[0] = listKey
	for idx, value := range list {
		args[idx*2+1] = idx
		args[idx*2+2] = value
	}
	if _, err = conn.Do("ZADD", args...); err != nil {
		err = fmt.Errorf("setStringList conn.Do(SADD, %s, %v) error(%v)", listKey, list, err)
		log.ErrLog("", err)
	}
	_, err = conn.Do("EXPIRE", listKey, expireTime)
	if err != nil {
		log.ErrLog("", err)
	}
	return err
}

//获取redis列表name id映射列表
func getNameList(conn redis.Conn, listKey string, expireTime int) ([]KeyValue, error) {
	rsp := make([]KeyValue, 0)
	ok, err := redis.Bool(conn.Do("EXPIRE", listKey, expireTime))
	if err != nil && err != redis.ErrNil {
		err = fmt.Errorf("getNameList conn.Do(EXPIRE, %s) error(%v)", listKey, err)
		log.ErrLog("", err)
		return nil, err
	}
	if !ok {
		// 不存在或者为空都会重新去DB获取
		err = redis.ErrNil
		return nil, err
	}
	list, err := redis.Strings(conn.Do("HGETALL", listKey))
	if err != nil && err != redis.ErrNil {
		if e, ok := err.(*redis.Error); ok && strings.Index(e.Error(), "WRONGTYPE") == -1 {
			err = fmt.Errorf("redisGetNameList conn.Do(HGETALL, %s) error(%v)", listKey, err)
			log.ErrLog("", err)
			return rsp, err
		}
		err = nil
		return rsp, err
	}
	size := len(list)
	if size%2 != 0 {
		err = errors.New("缓存内容数量错误")
		log.ErrLog("", err)
		return rsp, err
	}
	for i := 0; i+1 < size; i = i + 2 {
		rsp = append(rsp, KeyValue{Key: list[i], Value: list[i+1]})
	}
	return rsp, err
}

func getNameValue(conn redis.Conn, listKey string, key string, expireTime int) (*KeyValue, error) {
	rsp := &KeyValue{}
	rsp.Key = key
	ok, err := redis.Bool(conn.Do("EXPIRE", listKey, expireTime))
	if err != nil && err != redis.ErrNil {
		err = fmt.Errorf("getNameList conn.Do(EXPIRE, %s) error(%v)", listKey, err)
		log.ErrLog("", err)
		return rsp, err
	}
	if !ok {
		// 不存在或者为空都会重新去DB获取
		err = redis.ErrNil
		return rsp, err
	}
	rsp.Value, err = redis.String(conn.Do("HGET", listKey, key))
	if err != nil && err != redis.ErrNil {
		if e, ok := err.(*redis.Error); ok && strings.Index(e.Error(), "WRONGTYPE") == -1 {
			err = fmt.Errorf("redisGetNameList conn.Do(HGETALL, %s) error(%v)", listKey, err)
			log.ErrLog("", err)
			return rsp, err
		}
		err = nil
	}
	return rsp, err
}

//设置redis列表name id映射
func setNameList(conn redis.Conn, listKey string,
	list []KeyValue, expireTime int) (err error) {
	if len(list) <= 0 {
		// 设置哨兵
		if _, err = conn.Do("SETEX", listKey, expireTime, "emptylist"); err != nil {
			err = fmt.Errorf("redisSetNameList conn.Do(SETEX, %s, %v) error(%v)", listKey, list, err)
			log.ErrLog("", err)
		}
		return err
	}
	if _, err = conn.Do("DEL", listKey); err != nil && err != redis.ErrNil {
		err = fmt.Errorf("redisSetNameList conn.Do(DEL, %s) error(%v)", listKey, err)
		log.ErrLog("", err)
		return err
	}
	args := make([]interface{}, len(list)*2+1)
	args[0] = listKey
	for i, f := range list {
		args[i*2+1] = f.Key
		args[i*2+2] = f.Value
	}
	if _, err = conn.Do("HMSET", args...); err != nil {
		err = fmt.Errorf("redisSetNameList conn.Do(HMSET, %v) error(%v)", args, err)
		log.ErrLog("", err)
		return err
	}
	_, err = conn.Do("EXPIRE", listKey, expireTime)
	if err != nil {
		log.ErrLog("", err)
	}
	return err
}

func setKeyString(conn redis.Conn, key string,
	data string, expireTime int) error {
	_, err := conn.Do("SET", key, data, SetWithExpireTime, expireTime)
	if err != nil {
		err = fmt.Errorf("RedisSetKeyString conn.Do(SET, %v) error(%v)", key, err)
		log.ErrLog("", err)
	}
	return err
}

func getKeyString(conn redis.Conn, key string,
	expireTime int) (string, error) {
	var data string
	ok, err := redis.Bool(conn.Do("EXPIRE", key, expireTime))
	if err != nil && err != redis.ErrNil {
		err = fmt.Errorf("getKeyString conn.Do(EXPIRE, %s) error(%v)", key, err)
		log.ErrLog("", err)
		return data, err
	}
	if !ok {
		err = redis.ErrNil
		return data, err
	}

	data, err = redis.String(conn.Do("GET", key))
	if err != nil && err != redis.ErrNil {
		err = fmt.Errorf("getKeyString conn.Do(GET, %s) error(%v)", key, err)
		log.ErrLog("", err)
		return data, err
	}
	return data, nil
}

//10分钟人数统计
func rpushOnlineCount(conn redis.Conn, key string, count int64) (err error) {
	data, err := redis.Int64(conn.Do("LLEN", key))
	if err != nil && err != redis.ErrNil {
		err = fmt.Errorf("getKeyString conn.Do(LLEN, %s) error(%v)", key, err)
		log.ErrLog("", err)
		return err
	}
	//只需要前十分钟的数据，每一分钟进行存入
	if data >= 10 {

		forLen := data - 9
		for i := forLen; i > 0; i-- {
			_, err := redis.Bool(conn.Do("LPOP", key))
			if err != nil && err != redis.ErrNil {
				err = fmt.Errorf("getKeyString conn.Do(LPOP, %s) error(%v)", key, err)
				log.ErrLog("", err)
				return err
			}
		}
	}
	_, err = redis.Bool(conn.Do("RPUSH", key, count))
	if err != nil && err != redis.ErrNil {
		err = fmt.Errorf("RPUSH conn.Do(RPUSH, %s) error(%v)", key, err)
		log.ErrLog("", err)
		return err
	}

	//if !ok {
	//	err = fmt.Errorf("not found key(%s) ", key)
	//	return err
	//}
	return
}

//
func getLenOnlineCount(conn redis.Conn, key string) (data []int64, err error) {
	data, err = redis.Int64s(conn.Do("LRANGE", key, "0", "-1"))
	if err != nil && err != redis.ErrNil {
		if e, ok := err.(*redis.Error); ok && strings.Index(e.Error(), "WRONGTYPE") == -1 {
			err = fmt.Errorf("getLenOnlineCount conn.Do(LRANGE, %s) error(%v)", key, err)
			log.ErrLog("", err)
			return data, err
		}
	}
	return
}

//func getLenOnlineCount(conn redis.Conn, key string) (data []string,err error) {
//	list, err := redis.Strings(conn.Do("ZRANGE", key, "0", "-1"))
//	if err != nil && err != redis.ErrNil {
//		if e, ok := err.(*redis.Error); ok && strings.Index(e.Error(), "WRONGTYPE") == -1 {
//			err = fmt.Errorf("getStringList conn.Do(ZRANGE, %s) error(%v)", key, err)
//			log.ErrLog("", err)
//			return list, err
//		}
//	}
//	return
//}

func rpushList(conn redis.Conn, key string, count int64) (err error) {
	data, err := redis.Int64(conn.Do("LLEN", key))
	if err != nil && err != redis.ErrNil {
		err = fmt.Errorf("getKeyString conn.Do(LLEN, %s) error(%v)", key, err)
		log.ErrLog("", err)
		return err
	}
	//只需要前十分钟的数据，每一分钟进行存入
	if data >= 10 {
		forLen := data - 9
		for i := forLen; i == 0; i-- {
			ok, err := redis.Bool(conn.Do("LPOP", key))
			if err != nil && err != redis.ErrNil {
				err = fmt.Errorf("getKeyString conn.Do(LPOP, %s) error(%v)", key, err)
				log.ErrLog("", err)
				return err
			}
			if !ok {
				err = fmt.Errorf("not found key(%s) ", key)
				return err
			}
		}
	}
	ok, err := redis.Bool(conn.Do("RPUSH", key, count))
	if err != nil && err != redis.ErrNil {
		err = fmt.Errorf("RPUSH conn.Do(RPUSH, %s) error(%v)", key, err)
		log.ErrLog("", err)
		return err
	}
	if !ok {
		err = fmt.Errorf("not found key(%s) ", key)
		return err
	}
	return
}

func setSetID(conn redis.Conn, key string, id int64) (err error) {
	_, err = redis.Bool(conn.Do("SADD", key, id))
	if err != nil && err != redis.ErrNil {
		err = fmt.Errorf("SADD conn.Do(SADD, %s) error(%v)", key, err)
		log.ErrLog("", err)
		return err
	}
	//if !ok {
	//	err = fmt.Errorf("not found key(%s) ", key)
	//	return err
	//}
	return
}

func getSetCount(conn redis.Conn, key string) (count int64, err error) {
	//ok, err := redis.Bool(conn.Do("EXISTS", key))
	//if err != nil && err != redis.ErrNil {
	//	err = fmt.Errorf("SADD conn.Do(SADD, %s) error(%v)", key, err)
	//	log.ErrLog("", err)
	//	return count,err
	//}
	//if !ok  {
	//	//不存在或者获取失败，返回0
	//	err = fmt.Errorf("not found key(%s) ", key)
	//	return count, err
	//}
	count, err = redis.Int64(conn.Do("SCARD", key))
	if err != nil && err != redis.ErrNil {
		err = fmt.Errorf("SCARD conn.Do(SCARD, %s) error(%v)", key, err)
		log.ErrLog("", err)
		return count, err
	}
	return
}

//获取key
func getRegexpKey(conn redis.Conn, key string) ([]string, error) {
	keys := make([]string, 0)
	var err error
	keys, err = redis.Strings(conn.Do("Keys", key))
	if err != nil {
		err = fmt.Errorf("regexpDelKey:%v;key=%s", err, key)
		log.ErrLog("", err)
		return keys, err
	}
	return keys, nil
}

func hReset(conn redis.Conn, key string, field string, expireTime int) error {
	args := make([]interface{}, 3)
	args[0] = key
	args[1] = field
	args[2] = 0
	if _, err := conn.Do("HSET", args...); err != nil {
		err = fmt.Errorf("setKeyInfo conn.Do(HSET, %v) error(%v)", args, err)
		log.ErrLog("", err)
		return err
	}
	_, err := conn.Do("EXPIRE", key, expireTime)
	if err != nil {
		log.ErrLog("", err)
	}
	return err
}

func hIncrBy(conn redis.Conn, key string, field string, num int64, expireTime int) error {
	args := make([]interface{}, 3)
	args[0] = key
	args[1] = field
	args[2] = num
	if _, err := conn.Do("HINCRBY", args...); err != nil {
		err = fmt.Errorf("setKeyInfo conn.Do(HINCRBY, %v) error(%v)", args, err)
		log.ErrLog("", err)
		return err
	}
	_, err := conn.Do("EXPIRE", key, expireTime)
	if err != nil {
		log.ErrLog("", err)
	}
	return err
}

func hGetNum(conn redis.Conn, key string, field string, expireTime int) (int64, error) {
	var (
		rel int64
		ok  bool
		err error
	)
	ok, err = redis.Bool(conn.Do("EXPIRE", key, expireTime))
	if err != nil && err != redis.ErrNil {
		err = fmt.Errorf("hGetNum conn.Do(EXPIRE, %s) error(%v)", key, err)
		log.ErrLog("", err)
		return rel, err
	}
	if !ok {
		return rel, nil
	}
	args := make([]interface{}, 2)
	args[0] = key
	args[1] = field
	rel, err = redis.Int64(conn.Do("HMGET", args...))
	if err != nil && err != redis.ErrNil {
		if e, ok := err.(*redis.Error); ok && strings.Index(e.Error(), "WRONGTYPE") == -1 {
			err = fmt.Errorf("redisGetKeyInfo conn.Do(HMGET, %s,%v) error(%v)", key, field, err)
			log.ErrLog("", err)
			return rel, err
		}
		err = nil
		return rel, err
	}
	return rel, nil
}

func hDel(conn redis.Conn, key string, fields []string) error {
	args := make([]interface{}, len(fields)+1)
	args[0] = key
	for i := 1; i <= len(fields); i++ {
		args[i] = fields[i-1]
	}
	if _, err := conn.Do("HDEL", args...); err != nil {
		err = fmt.Errorf("delKey:%v;key=%s", err, key)
		log.ErrLog("", err)
		return err
	}
	return nil
}
