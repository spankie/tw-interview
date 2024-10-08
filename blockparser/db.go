package blockparser

import (
	"errors"
	"strings"
	"sync"
)

var ErrInvalidKey = errors.New("invalid key")

type memoryStore[T any] struct {
	mu sync.RWMutex

	// NOTE: sync.Map is not used because from the documentation,
	// this usecase does not require it.
	data map[string][]T
}

func newMemoryDataStore[T any]() *memoryStore[T] {
	return &memoryStore[T]{
		data: make(map[string][]T),
	}
}

func (s *memoryStore[T]) Add(key string, value []T) error {
	if strings.TrimSpace(key) == "" {
		return ErrInvalidKey
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.data[key]; !ok {
		s.data[key] = value
		return nil
	}

	s.data[key] = append(s.data[key], value...)

	return nil
}

func (s *memoryStore[T]) Get(key string) ([]T, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	item, ok := s.data[key]

	return item, ok
}

func (s *memoryStore[T]) GetKeys() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	keys := make([]string, 0, len(s.data))

	for k := range s.data {
		keys = append(keys, k)
	}

	return keys
}
