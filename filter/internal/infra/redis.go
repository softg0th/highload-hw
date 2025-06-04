package infra

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

type Redis struct {
	rdb *redis.Client
}

func NewRedisClient(host string, port string) *Redis {
	redisAddr := fmt.Sprintf("%s:%s", host, port)
	rdb := redis.NewClient(&redis.Options{
		Addr: redisAddr,
		DB:   0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		fmt.Printf("Redis connection error: %v\n", err)
	}

	return &Redis{rdb: rdb}
}

func (r *Redis) Get(key string) (string, error) {
	ctx := context.Background()
	return r.rdb.Get(ctx, key).Result()
}

func (r *Redis) Set(key, value string, ttl time.Duration) error {
	ctx := context.Background()
	return r.rdb.Set(ctx, key, value, ttl).Err()
}
