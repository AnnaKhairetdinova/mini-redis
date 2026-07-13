package store

import (
	"sync"
	"time"
)

type entry struct {
	value     string
	expiresAt time.Time
}

type Store struct {
	mu   sync.RWMutex
	data map[string]entry
}

func New() *Store {
	return &Store{
		sync.RWMutex{},
		make(map[string]entry),
	}
}

func (s *Store) Set(key, value string, ttl time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()
	var expires time.Time

	if ttl > 0 {
		expires = time.Now().Add(ttl)
	}

	s.data[key] = entry{value, expires}
}

func (s *Store) Get(key string) (string, bool) {
	//ключ существует? не протух? Если протух — вернуть ("", false)
	s.mu.RLock()
	defer s.mu.RUnlock()

	if _, ok := s.data[key]; !ok {
		return "", false
	}

	if !s.data[key].expiresAt.IsZero() && time.Now().After(s.data[key].expiresAt) {
		return "", false
	}

	return key, true
}

func (s *Store) Del(key string) bool {
	s.mu.Lock()

	if _, ok := s.data[key]; ok {
		delete(s.data, key)
		return true
	}

	s.mu.Unlock()

	return false
}

func (s *Store) Keys() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	res := make([]string, 0, len(s.data))
	for k := range s.data {
		res = append(res, k)
	}

	return res
}
