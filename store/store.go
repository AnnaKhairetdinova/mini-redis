package store

import (
	"context"
	"sync"
	"sync/atomic"
	"time"
)

type entry struct {
	value     string
	expiresAt time.Time
}

type Store struct {
	mu      sync.RWMutex
	data    map[string]entry
	Clients atomic.Int64
}

func New() *Store {
	return &Store{
		data: make(map[string]entry),
	}
}

func (s *Store) Set(key string, value string, ttl time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()
	var expires time.Time

	if ttl > 0 {
		expires = time.Now().Add(ttl)
	} else if ttl < 0 {
		expires = time.Now().Add(-1 * time.Nanosecond)
	}

	s.data[key] = entry{value, expires}
}

func (s *Store) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if _, ok := s.data[key]; !ok {
		return "", false
	}

	if !s.data[key].expiresAt.IsZero() && time.Now().After(s.data[key].expiresAt) {
		return "", false
	}

	return s.data[key].value, true
}

func (s *Store) Del(key string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.data[key]; ok {
		delete(s.data, key)
		return true
	}

	return false
}

func (s *Store) Keys() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	res := make([]string, 0, len(s.data))
	for k, e := range s.data {
		if e.expiresAt.IsZero() || time.Now().Before(e.expiresAt) {
			res = append(res, k)
		}
	}

	return res
}

func (s *Store) StartCleaner(ctx context.Context, interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				s.mu.Lock()

				for key, val := range s.data {
					if !val.expiresAt.IsZero() && time.Now().After(val.expiresAt) {
						delete(s.data, key)
					}
				}
				s.mu.Unlock()

			case <-ctx.Done():
				return
			}
		}
	}()
}
