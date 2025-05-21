package cache

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
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

func (c *Cache) Redis(ctx context.Context) *redis.Client {
	return c.rdb
}

func NewRedis(conf *viper.Viper) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:         conf.GetString("data.redis.addr"),
		Password:     conf.GetString("data.redis.password"),
		DB:           conf.GetInt("data.redis.db"),
		ReadTimeout:  time.Duration(conf.GetInt("data.redis.read_timeout")) * time.Millisecond,
		WriteTimeout: time.Duration(conf.GetInt("data.redis.write_timeout")) * time.Millisecond,
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
