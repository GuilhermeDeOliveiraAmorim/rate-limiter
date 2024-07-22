package main

import (
	"time"
)

type RateLimiterStorage interface {
	Increment(key string, duration time.Duration) (int64, error)
}

type RateLimiter struct {
	storage            RateLimiterStorage
	rateLimitIP        int
	rateLimitToken     int
	blockDurationIP    time.Duration
	blockDurationToken time.Duration
}

func NewRateLimiter(storage RateLimiterStorage, rateLimitIP, rateLimitToken int, blockDurationIP, blockDurationToken time.Duration) *RateLimiter {
	return &RateLimiter{
		storage:            storage,
		rateLimitIP:        rateLimitIP,
		rateLimitToken:     rateLimitToken,
		blockDurationIP:    blockDurationIP,
		blockDurationToken: blockDurationToken,
	}
}

func (rl *RateLimiter) AllowAccessByIP(ip string) bool {
	key := "rate_limiter_ip:" + ip

	count, err := rl.storage.Increment(key, rl.blockDurationIP)
	if err != nil {
		return false
	}

	if count > int64(rl.rateLimitIP) {
		return false
	}

	return true
}

func (rl *RateLimiter) AllowAccessByToken(token string) bool {
	key := "rate_limiter_token:" + token

	count, err := rl.storage.Increment(key, rl.blockDurationToken)
	if err != nil {
		return false
	}

	if count > int64(rl.rateLimitToken) {
		return false
	}

	return true
}
