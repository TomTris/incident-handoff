package main

import (
	"context"
	"strconv"
)

type MemoryUserStore struct {
	users     map[string]User // username - User
	currentID int
}

func NewMemoryUserStoreWithSeed(seed []User) *MemoryUserStore {
	m := make(map[string]User, len(seed))
	for _, u := range seed {
		m[u.Username] = u
	}
	return &MemoryUserStore{users: m}
}

func NewMemoryUserStore() *MemoryUserStore {
	m := make(map[string]User)
	return &MemoryUserStore{users: m}
}

func (s *MemoryUserStore) Create(ctx context.Context, u User) (User, error) {
	_, ok := s.users[u.Username]
	if ok == true {
		return User{}, ErrUserAlreadyExist
	}

	s.currentID++
	ID := UserPrefix + strconv.Itoa(s.currentID)
	u.ID = ID
	s.users[u.Username] = u
	return u, nil
}

func (s *MemoryUserStore) GetByUsername(_ context.Context, username string) (User, error) {
	u, ok := s.users[username]
	if ok == false {
		return User{}, ErrUserNotFound
	}
	return u, nil
}
