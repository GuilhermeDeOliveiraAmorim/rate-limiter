package main

import (
	"sync"
	"time"
)

type MemoryStorage struct {
	data map[string]int64
	mu   sync.Mutex
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		data: make(map[string]int64),
	}
}

func (ms *MemoryStorage) Increment(key string, duration time.Duration) (int64, error) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	ms.data[key]++
	count := ms.data[key]

	if count == 1 {
		go func() {
			time.Sleep(duration)
			ms.mu.Lock()
			delete(ms.data, key)
			ms.mu.Unlock()
		}()
	}

	return count, nil
}
