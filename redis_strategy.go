package main

import (
	"time"

	"github.com/go-redis/redis/v8"
	"golang.org/x/net/context"
)

type RedisStorage struct {
	client *redis.Client
}

func NewRedisStorage(host, port string) *RedisStorage {
	addr := host + ":" + port
	client := redis.NewClient(&redis.Options{
		Addr:        addr,
		DB:          0,
		DialTimeout: 10 * time.Millisecond,
	})
	return &RedisStorage{client: client}
}

func (rs *RedisStorage) Increment(key string, duration time.Duration) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	count, err := rs.client.Incr(ctx, key).Result()
	if err != nil {
		return 0, err
	}

	if count == 1 {
		rs.client.Expire(ctx, key, duration)
	}

	return count, nil
}
