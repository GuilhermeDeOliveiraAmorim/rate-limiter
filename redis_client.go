package main

import (
	"github.com/go-redis/redis/v8"
)

func NewRedisClient(host, port string) *redis.Client {
	addr := host + ":" + port
	return redis.NewClient(&redis.Options{
		Addr: addr,
		DB:   0, // Use default DB
	})
}
