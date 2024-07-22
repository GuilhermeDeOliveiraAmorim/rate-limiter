package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func TestRateLimiterWithRedis(t *testing.T) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:        "localhost:6379",
		DB:          0,
		DialTimeout: 10 * time.Millisecond,
	})
	ctx := context.Background()
	redisClient.FlushDB(ctx)

	storage := &RedisStorage{client: redisClient}
	rateLimiter := NewRateLimiter(storage, 2, 5, 10*time.Second, 10*time.Second)

	handler := rateLimiter.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	}))

	server := httptest.NewServer(handler)
	defer server.Close()

	client := &http.Client{}

	// Teste limite por IP
	for i := 0; i < 3; i++ {
		req, _ := http.NewRequest("GET", server.URL, nil)
		req.RemoteAddr = "192.168.1.1:12345"
		resp, _ := client.Do(req)

		if i < 2 {
			assert.Equal(t, http.StatusOK, resp.StatusCode)
		} else {
			assert.Equal(t, http.StatusTooManyRequests, resp.StatusCode)
		}
	}

	// Teste limite por Token
	for i := 0; i < 6; i++ {
		req, _ := http.NewRequest("GET", server.URL, nil)
		req.Header.Set("API_KEY", "token123")
		resp, _ := client.Do(req)

		if i < 5 {
			assert.Equal(t, http.StatusOK, resp.StatusCode)
		} else {
			assert.Equal(t, http.StatusTooManyRequests, resp.StatusCode)
		}
	}
}

func TestRateLimiterWithMemory(t *testing.T) {
	storage := NewMemoryStorage()
	rateLimiter := NewRateLimiter(storage, 2, 5, 10*time.Second, 10*time.Second)

	handler := rateLimiter.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	}))

	server := httptest.NewServer(handler)
	defer server.Close()

	client := &http.Client{}

	// Teste limite por IP
	for i := 0; i < 3; i++ {
		req, _ := http.NewRequest("GET", server.URL, nil)
		req.RemoteAddr = "192.168.1.1:12345"
		resp, _ := client.Do(req)

		if i < 2 {
			assert.Equal(t, http.StatusOK, resp.StatusCode)
		} else {
			assert.Equal(t, http.StatusTooManyRequests, resp.StatusCode)
		}
	}

	// Teste limite por Token
	for i := 0; i < 6; i++ {
		req, _ := http.NewRequest("GET", server.URL, nil)
		req.Header.Set("API_KEY", "token123")
		resp, _ := client.Do(req)

		if i < 5 {
			assert.Equal(t, http.StatusOK, resp.StatusCode)
		} else {
			assert.Equal(t, http.StatusTooManyRequests, resp.StatusCode)
		}
	}
}
