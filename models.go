package main

import (
	"sync"
	"time"
)

// User represents a user in our system
type User struct {
	ID        int       `toon:"id"`
	Name      string    `toon:"name"`
	Email     string    `toon:"email"`
	CreatedAt time.Time `toon:"createdAt"`
}

// UserStore manages users in memory
type UserStore struct {
	mu      sync.RWMutex
	users   map[int]*User
	nextID  int
}

// NewUserStore creates a new user store
func NewUserStore() *UserStore {
	return &UserStore{
		users:  make(map[int]*User),
		nextID: 1,
	}
}

// GetAll returns all users
func (s *UserStore) GetAll() []*User {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	users := make([]*User, 0, len(s.users))
	for _, user := range s.users {
		users = append(users, user)
	}
	return users
}

// Get returns a user by ID
func (s *UserStore) Get(id int) (*User, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	user, exists := s.users[id]
	return user, exists
}

// Create adds a new user
func (s *UserStore) Create(user *User) *User {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	user.ID = s.nextID
	user.CreatedAt = time.Now()
	s.users[s.nextID] = user
	s.nextID++
	
	return user
}

// Update modifies an existing user
func (s *UserStore) Update(id int, user *User) (*User, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	if _, exists := s.users[id]; !exists {
		return nil, false
	}
	
	user.ID = id
	s.users[id] = user
	return user, true
}

// Delete removes a user
func (s *UserStore) Delete(id int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	if _, exists := s.users[id]; !exists {
		return false
	}
	
	delete(s.users, id)
	return true
}
