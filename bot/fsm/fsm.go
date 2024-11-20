package fsm

import "sync"

type Session struct {
	users map[int64]*UserState
	mu    sync.Mutex
}

type UserState struct {
	State string
}

func NewSession() *Session {
	return &Session{
		users: make(map[int64]*UserState),
	}
}

func (s *Session) GetState(userID int64) *UserState {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.users[userID]; !exists {
		s.users[userID] = &UserState{State: "initial"}
	}
	return s.users[userID]
}

func (s *Session) SetState(userID int64, state string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.users[userID]; !exists {
		s.users[userID].State = state
	}
	s.users[userID].State = state
}