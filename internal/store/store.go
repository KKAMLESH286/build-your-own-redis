//Store keys, values, and expiration information

package store

import (
	"sync"
	"time"
)

type cacheItem struct {
	value     string
	expiresAt *time.Time
}

type Store struct {
	mu   sync.RWMutex
	data map[string]cacheItem
}

func NewStore() *Store {
	s := &Store{
		data: make(map[string]cacheItem),
	}

	go s.startExpirationRoutine()

	return s

}

func (s *Store) Set(key, value string, duration *time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var expiresAt *time.Time
	if duration != nil {
		t := time.Now().Add(*duration)
		expiresAt = &t
	}

	s.data[key] = cacheItem{
		value:     value,
		expiresAt: expiresAt,
	}
}

func (s *Store) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	item, exists := s.data[key]
	if !exists {
		return "", false
	}
	if item.expiresAt != nil && time.Now().After(*item.expiresAt) {
		delete(s.data, key)
		return "", false
	}
	return item.value, true
}

func (s *Store) RemoveExpired() {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	for key, item := range s.data {
		if item.expiresAt != nil && now.After(*item.expiresAt) {
			delete(s.data, key)
		}
	}
}

func (s *Store) startExpirationRoutine() {
	ticker := time.NewTicker(1 * time.Second)
	for range ticker.C {
		s.RemoveExpired()
	}
}

func (s *Store) Delete(keys []string) int {
	s.mu.Lock()
	defer s.mu.Unlock()

	count := 0
	for _, key := range keys {
		if _, exists := s.data[key]; exists {
			delete(s.data, key)
			count++
		}
	}
	return count
}
