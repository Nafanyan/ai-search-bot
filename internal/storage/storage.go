package storage

import "sync"

type ChatState int

const (
	StateIdle ChatState = iota
	StateSearching
)

type ChatStorage interface {
	GetState(chatID int64) ChatState
	SetState(chatID int64, state ChatState)
}

type InMemoryChatStorage struct {
	mu     sync.Mutex
	states map[int64]ChatState
}

func NewInMemoryChatStorage() *InMemoryChatStorage {
	return &InMemoryChatStorage{
		states: make(map[int64]ChatState),
	}
}

func (s *InMemoryChatStorage) GetState(chatID int64) ChatState {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.states[chatID]
}

func (s *InMemoryChatStorage) SetState(chatID int64, state ChatState) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.states[chatID] = state
}
