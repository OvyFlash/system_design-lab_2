package redis

import (
	"fmt"
	"lab_2/config"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/sheophe/go-rejson"
	"github.com/sheophe/redisearch-go/redisearch"
)

type Redis struct {
	Pool   *redis.Pool
	Search *redisearch.Client
	Json   *rejson.Handler
}

func NewRedis(conf config.RedisDatabase) (r *Redis, err error) {
	r = &Redis{}
	r.Pool = &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial(
				"tcp",
				fmt.Sprintf("%s:%d", conf.Host, conf.Port),
				redis.DialDatabase(conf.DB),
				redis.DialPassword(conf.Password),
			)
			if err != nil {
				return nil, err
			}
			r.Search = redisearch.NewClient(c)
			r.Json = rejson.NewReJSONHandler()
			r.Json.SetRedigoClient(c)
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
	return
}

func (r *Redis) Set(key string, value any) (err error) {
	conn := r.Pool.Get()
	defer conn.Close()
	_, err = conn.Do("SET", key, value)
	if err != nil {
		return
	}
	return
}

func (r *Redis) SetWithExpiration(key string, value any, expiration time.Duration) (err error) {
	conn := r.Pool.Get()
	defer conn.Close()
	_, err = redis.String(conn.Do(
		"SET", key, value,
		"EX", expiration.Seconds(),
	))
	return
}

func (r *Redis) AddToSet(set, key string) (err error) {
	conn := r.Pool.Get()
	defer conn.Close()
	_, err = conn.Do("SADD", set, key)
	return
}

func (r *Redis) SetMap(key string, fields map[string]string) (err error) {
	conn := r.Pool.Get()
	defer conn.Close()
	var parsedFields redis.Args
	for k, v := range fields {
		parsedFields.Add(k).Add(v)
	}
	_, err = conn.Do("HSET", key, parsedFields)
	return
}

func (r *Redis) SetMapField(key, field, value string) (err error) {
	conn := r.Pool.Get()
	defer conn.Close()
	_, err = conn.Do("HSET", key, field, value)
	return
}

func (r *Redis) GetString(key string) (value string, err error) {
	conn := r.Pool.Get()
	defer conn.Close()
	return redis.String(conn.Do("GET", key))
}

func (r *Redis) GetInt(key string) (value int64, err error) {
	conn := r.Pool.Get()
	defer conn.Close()
	return redis.Int64(conn.Do("GET", key))
}

func (r *Redis) IsInSet(set, key string) (bool, error) {
	conn := r.Pool.Get()
	defer conn.Close()
	return redis.Bool(conn.Do("SISMEMBER", set, key))
}

func (r *Redis) GetMap(key string) (value map[string]string, err error) {
	conn := r.Pool.Get()
	defer conn.Close()
	return redis.StringMap(conn.Do("HGETALL", key))
}

func (r *Redis) GetMapField(key, field string) (value string, err error) {
	conn := r.Pool.Get()
	defer conn.Close()
	return redis.String(conn.Do("HGET", key, field))
}

func (r *Redis) Delete(key string) (err error) {
	conn := r.Pool.Get()
	defer conn.Close()
	_, err = conn.Do("DEL", key)
	return
}

func (r *Redis) DeleteFromSet(set, key string) (err error) {
	conn := r.Pool.Get()
	defer conn.Close()
	_, err = conn.Do("SREM", set, key)
	return
}

func (r *Redis) DeleteFromMap(key, field string) (err error) {
	conn := r.Pool.Get()
	defer conn.Close()
	_, err = conn.Do("HDEL", key, field)
	return
}

func (r *Redis) Keys(filter string) (keys []string, err error) {
	conn := r.Pool.Get()
	defer conn.Close()
	return redis.Strings(conn.Do("KEYS", filter))
}

func (r *Redis) KeysMap(filter string) (keys []string, err error) {
	conn := r.Pool.Get()
	defer conn.Close()
	return redis.Strings(conn.Do("HKEYS", filter))
}
