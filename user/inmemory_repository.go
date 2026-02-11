package user

import (
	"context"
	"fmt"
	"sync"
)

type InMemoryUserRepository struct {
	mu       sync.Mutex
	users    map[string]User
	auditLog map[string][]string
}

func NewInMemoryUserRepository(seed []User) *InMemoryUserRepository {
	users := make(map[string]User, len(seed))
	for _, u := range seed {
		users[u.ID] = u
	}
	return &InMemoryUserRepository{
		users:    users,
		auditLog: make(map[string][]string),
	}
}

func (r *InMemoryUserRepository) Exists(_ context.Context, id string) (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, ok := r.users[id]
	return ok, nil
}

func (r *InMemoryUserRepository) FindByID(_ context.Context, id string) (User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	u, ok := r.users[id]
	if !ok {
		return User{}, ErrUserNotFound
	}
	return u, nil
}

func (r *InMemoryUserRepository) IncrementFailedLogin(_ context.Context, id string) (User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	u, ok := r.users[id]
	if !ok {
		return User{}, ErrUserNotFound
	}
	u.FailedLoginCount++
	if u.FailedLoginCount >= maxFailedLoginCount {
		u.Locked = true
	}
	r.users[id] = u
	return u, nil
}

func (r *InMemoryUserRepository) ResetFailedLogin(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	u, ok := r.users[id]
	if !ok {
		return ErrUserNotFound
	}
	u.FailedLoginCount = 0
	r.users[id] = u
	return nil
}

func (r *InMemoryUserRepository) AppendAuditLog(_ context.Context, id string, result string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.users[id]; !ok {
		return fmt.Errorf("append audit log: %w", ErrUserNotFound)
	}
	r.auditLog[id] = append(r.auditLog[id], result)
	return nil
}

func (r *InMemoryUserRepository) AuditLog(id string) []string {
	r.mu.Lock()
	defer r.mu.Unlock()
	logs := r.auditLog[id]
	copied := make([]string, len(logs))
	copy(copied, logs)
	return copied
}
