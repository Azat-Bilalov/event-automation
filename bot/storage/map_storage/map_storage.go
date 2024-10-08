package storage

import (
	"event-automation/bot/lib/e"
	"fmt"
)

type Storage struct {
	mapping map[string]string
}

func New() *Storage {
	return &Storage{
		mapping: make(map[string]string),
	}
}

func (s *Storage) Save(key string, value string) error {
	s.mapping[key] = value
	return nil
}

func (s *Storage) IsExist(key string) (bool, error) {
	_, ok := s.mapping[key]
	return ok, nil
}

func (s *Storage) Delete(key string) error {
	if ok, _ := s.IsExist(key); !ok {
		return e.Wrap("Key not found", fmt.Errorf("key: %s", key))
	}
	delete(s.mapping, key)
	return nil
}

func (s *Storage) Get(key string) (string, error) {
	if ok, _ := s.IsExist(key); !ok {
		return "", e.Wrap("Key not found", fmt.Errorf("key: %s", key))
	}
	return s.mapping[key], nil
}
