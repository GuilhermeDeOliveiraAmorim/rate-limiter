package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
)

func TestRateLimiter(t *testing.T) {
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})

	rateLimiter := NewRateLimiter(redisClient, 2, 5, 10*time.Second, 10*time.Second)

	handler := rateLimiter.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	}))

	server := httptest.NewServer(handler)
	defer server.Close()

	client := &http.Client{}

	// Teste limite por IP
	req, _ := http.NewRequest("GET", server.URL, nil)
	req.RemoteAddr = "192.168.1.1:12345"

	resp, _ := client.Do(req)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	resp, _ = client.Do(req)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	resp, _ = client.Do(req)
	assert.Equal(t, http.StatusTooManyRequests, resp.StatusCode)

	// Teste limite por Token
	req, _ = http.NewRequest("GET", server.URL, nil)
	req.Header.Set("API_KEY", "token123")

	resp, _ = client.Do(req)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	resp, _ = client.Do(req)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	resp, _ = client.Do(req)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	resp, _ = client.Do(req)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	resp, _ = client.Do(req)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	resp, _ = client.Do(req)
	assert.Equal(t, http.StatusTooManyRequests, resp.StatusCode)
}
