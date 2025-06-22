package main

import (
	"sync"
)

type UserStore struct {
	mu    sync.RWMutex
	users map[string]*TOTP
}

func NewUserStore() *UserStore {
	return &UserStore{
		users: make(map[string]*TOTP),
	}
}

func (us *UserStore) GetUser(email string) (*TOTP, bool) {
	us.mu.RLock()
	defer us.mu.RUnlock()
	user, exists := us.users[email]
	return user, exists
}

func (us *UserStore) CreateUser(email string) (*TOTP, error) {
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