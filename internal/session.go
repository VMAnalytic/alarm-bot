package app

import "context"

type SessionState uint8

const (
	_ = iota

	AddingContacts
	RemovingContacts
)

type Session struct {
	ID    int
	State SessionState
}

func NewSession(ID int) *Session {
	return &Session{ID: ID}
}

func (s *Session) Transition(state SessionState) {
	s.State = state
}

type SessionStorage interface {
	Add(ctx context.Context, session *Session)
	ExistInState(ctx context.Context, ID int, state SessionState) bool
	Delete(ctx context.Context, ID int)
}
