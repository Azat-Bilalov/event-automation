package storage

import "sync"

type Storage interface {
	SetEmail(userID int64, email string)
	GetEmail(userID int64) string
	IsExist(userID int64) bool
}

type store struct {
	mu     sync.RWMutex
	emails map[int64]string
}

func NewStore() Storage {
	return &store{
		emails: make(map[int64]string),
	}
}

func (s *store) SetEmail(userID int64, email string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.emails[userID] = email
}

func (s *store) GetEmail(userID int64) string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	email := s.emails[userID]
	return email
}

func (s *store) IsExist(userID int64) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if _, ok := s.emails[userID]; ok {
		return true
	}
	return false
}
