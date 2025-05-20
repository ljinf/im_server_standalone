package cache

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

const (
	DefaultRedisDB               = 0
	DefaultRedisPwd              = ""
	DefaultRedisReadTimeout  int = 500
	DefaultRedisWriteTimeout int = 500
)

type Cache struct {
	rdb *redis.Client
}

func NewCache(r *redis.Client) *Cache {
	return &Cache{
		rdb: r,
	}
}

func NewRedis(addr string, opts ...RedisOption) *redis.Client {

	conf := &redisConf{
		pwd:          DefaultRedisPwd,
		db:           DefaultRedisDB,
		readTimeout:  DefaultRedisReadTimeout,
		writeTimeout: DefaultRedisWriteTimeout,
	}
	for _, v := range opts {
		v(conf)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     conf.pwd,
		DB:           conf.db,
		ReadTimeout:  time.Duration(conf.readTimeout) * time.Millisecond,
		WriteTimeout: time.Duration(conf.writeTimeout) * time.Millisecond,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		panic(fmt.Sprintf("redis error: %s", err.Error()))
	}

	return rdb
}

type redisConf struct {
	pwd          string
	db           int
	readTimeout  int
	writeTimeout int
}

type RedisOption func(r *redisConf)

func WithRedisPwd(pwd string) RedisOption {
	return func(r *redisConf) {
		r.pwd = pwd
	}
}

func WithRedisDB(db int) RedisOption {
	return func(r *redisConf) {
		r.db = db
	}
}

func WithRedisReadTimeout(timeout int) RedisOption {
	return func(r *redisConf) {
		r.readTimeout = timeout
	}
}

func WithRedisWriteTimeout(timeout int) RedisOption {
	return func(r *redisConf) {
		r.writeTimeout = timeout
	}
}
