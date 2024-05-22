package main

import (
	"time"

	"github.com/go-redis/redis/v8"
	"golang.org/x/net/context"
)

type RateLimiter struct {
	redisClient        *redis.Client
	rateLimitIP        int
	rateLimitToken     int
	blockDurationIP    time.Duration
	blockDurationToken time.Duration
}

func NewRateLimiter(redisClient *redis.Client, rateLimitIP, rateLimitToken int, blockDurationIP, blockDurationToken time.Duration) *RateLimiter {
	return &RateLimiter{
		redisClient:        redisClient,
		rateLimitIP:        rateLimitIP,
		rateLimitToken:     rateLimitToken,
		blockDurationIP:    blockDurationIP,
		blockDurationToken: blockDurationToken,
	}
}

func (rl *RateLimiter) AllowAccessByIP(ip string) bool {
	ctx := context.Background()
	key := "rate_limiter_ip:" + ip

	// Incrementa o contador de requisições
	count, err := rl.redisClient.Incr(ctx, key).Result()
	if err != nil {
		return false
	}

	// Se for a primeira requisição, seta a expiração
	if count == 1 {
		rl.redisClient.Expire(ctx, key, rl.blockDurationIP)
	}

	// Verifica se o número de requisições excede o limite
	if count > int64(rl.rateLimitIP) {
		return false
	}

	return true
}

func (rl *RateLimiter) AllowAccessByToken(token string) bool {
	ctx := context.Background()
	key := "rate_limiter_token:" + token

	// Incrementa o contador de requisições
	count, err := rl.redisClient.Incr(ctx, key).Result()
	if err != nil {
		return false
	}

	// Se for a primeira requisição, seta a expiração
	if count == 1 {
		rl.redisClient.Expire(ctx, key, rl.blockDurationToken)
	}

	// Verifica se o número de requisições excede o limite
	if count > int64(rl.rateLimitToken) {
		return false
	}

	return true
}
