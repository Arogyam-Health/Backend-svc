package token

import (
	"sync"
	"time"
)

type Token struct {
	AccessToken string
	ExpiresAt   time.Time
}

type TokenRuntime struct {
	mu    sync.RWMutex
	token Token
}

func NewRuntime() *TokenRuntime {
	return &TokenRuntime{}
}

func (s *TokenRuntime) Get() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.token.AccessToken
}

func (s *TokenRuntime) Set(token Token) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.token = token
}

func (s *TokenRuntime) IsValid() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.token.AccessToken == "" {
		return false
	}
	return time.Until(s.token.ExpiresAt) < 7*24*time.Hour
}
