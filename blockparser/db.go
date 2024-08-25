package blockparser

import (
	"fmt"
	"log"
	"strings"
	"sync"
)

var InvalidKeyError = fmt.Errorf("invalid key")

// TODO: make the value type an interface that can be validated
type memoryStore[T any] struct {
	mu sync.RWMutex
	// NOTE: sync.Map is not used because....
	data map[string][]T
}

func NewMemoryDataStore[T any]() *memoryStore[T] {
	return &memoryStore[T]{
		data: make(map[string][]T),
	}
}

func (s *memoryStore[T]) Add(key string, value []T) error {
	if strings.TrimSpace(key) == "" {
		return InvalidKeyError
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

// TODO: remove this or rename it to something useful
func (s *memoryStore[T]) printData() {
	for i, v := range s.data {
		log.Printf("addr %s:\n----------\n%v\n----------", i, len(v))
	}
}
