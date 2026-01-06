package redis

import (
	"context"
	"log"

	"github.com/dermot10/code-reviewer/backend_go/config"
	"github.com/redis/go-redis/v9"
)

const (
	CachePrefix = "cache:"
	QueuePrefix = "queue:"
)

type RedisClient struct {
	Rdb *redis.Client
}

func NewRedisService(cfg *config.Config) (*RedisClient, error) {
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

func (r *RedisClient) SetCache(ctx context.Context, key string, value []byte) error {
	// type:item - value
	return r.Rdb.Set(ctx, CachePrefix+key, value, 0).Err()
}

func (r *RedisClient) GetCache(ctx context.Context, key string) (string, error) {
	return r.Rdb.Get(ctx, CachePrefix+key).Result()
}

func (r *RedisClient) DelKey(ctx context.Context, key string) error {
	return r.Rdb.Del(ctx, key).Err()
}
func (r *RedisClient) PushQueue(ctx context.Context, data []byte) error {
	return r.Rdb.LPush(ctx, QueuePrefix+"tasks", data).Err()
}

func (r *RedisClient) PopQueue(ctx context.Context) (string, error) {
	return r.Rdb.RPop(ctx, QueuePrefix+"tasks").Result()
}
