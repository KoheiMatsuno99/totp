package main

import (
	"sync"
)

type UserStore struct {
	mu    sync.RWMutex
	users map[Email]*TOTP
}

func NewUserStore() *UserStore {
	return &UserStore{
		users: make(map[Email]*TOTP),
	}
}

func (us *UserStore) GetUser(email Email) (*TOTP, bool) {
	us.mu.RLock()
	defer us.mu.RUnlock()
	user, exists := us.users[email]
	return user, exists
}

func (us *UserStore) CreateUser(email Email) (*TOTP, error) {
	us.mu.Lock()
	defer us.mu.Unlock()

	secret, err := GenerateSecret()
	if err != nil {
		return nil, err
	}

	totp := NewTOTP(secret)
	us.users[email] = totp
	return totp, nil
}

type Server struct {
	userStore *UserStore
}

func NewServer() *Server {
	return &Server{
		userStore: NewUserStore(),
	}
}