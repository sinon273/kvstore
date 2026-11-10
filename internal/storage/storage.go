package storage

import (
	"context"
	"sync"
	"time"
)

type Store struct {
	mu      sync.RWMutex
	data    map[string]string
	expires map[string]time.Time
}

func New() *Store {
	return &Store{
		data:    make(map[string]string),
		expires: make(map[string]time.Time),
	}
}

func (s *Store) Set(key, value string, ttl time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = value

	if ttl > 0 {
		s.expires[key] = time.Now().Add(ttl)
	} else {
		delete(s.expires, key)
	}
}

func (s *Store) Get(key string) (string, bool) {
	s.mu.RLock()
	exp, exists := s.expires[key]
	s.mu.RUnlock()

	if exists && time.Now().After(exp) {
		s.mu.Lock()
		delete(s.data, key)
		delete(s.expires, key)
		s.mu.Unlock()
		return "", false
	}

	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.data[key]
	return v, ok
}

func (s *Store) Delete(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data, key)
	delete(s.expires, key)
}

func (s *Store) StartCleanup(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				s.mu.Lock()
				now := time.Now()
				for key, exp := range s.expires {
					if now.After(exp) {
						delete(s.data, key)
						delete(s.expires, key)
					}
				}
				s.mu.Unlock()
			}
		}
	}()
}

func (s *Store) TTL(key string) int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if _, ok := s.data[key]; !ok {
		return -2
	}
	exp, ok := s.expires[key]
	if !ok {
		return -1
	}

	return int(time.Until(exp).Seconds())
}
