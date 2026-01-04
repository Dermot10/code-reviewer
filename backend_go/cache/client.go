package cache

import (
	"context"
	"log"

	"github.com/dermot10/code-reviewer/backend_go/config"
	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	Rdb *redis.Client
}

func NewCacheService(cfg *config.Config) (*RedisClient, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddr,
	})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}
	log.Print("Successfully connected to redis")
	return &RedisClient{Rdb: rdb}, nil
}

func (r *RedisClient) Close() error {
	return r.Rdb.Close()
}

// TODO - reduce redundant code in auth service but adding setJSON and getJSON to redis client
