package inmemory

import (
	"fmt"
	"sync"

	"github.com/romreign/PR-Reviewer-Assignment-Service/internal/api"
)

type UserRepository struct {
	mu    sync.RWMutex
	users map[string]*api.User
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		users: make(map[string]*api.User),
	}
}

func (r *UserRepository) FindUserByID(userID string) (*api.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, ok := r.users[userID]
	if !ok {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}

func (r *UserRepository) UpdateUserStatus(userID string, status bool) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	user, ok := r.users[userID]
	if !ok {
		return fmt.Errorf("user not found")
	}
	user.IsActive = status
	return nil
}

func (r *UserRepository) AddUser(user *api.User) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.users[user.UserId] = user
}

func (r *UserRepository) GetAllUsers() ([]api.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []api.User
	for _, user := range r.users {
		result = append(result, *user)
	}
	return result, nil
}
