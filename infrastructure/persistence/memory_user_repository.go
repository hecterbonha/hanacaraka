package persistence

import (
	"errors"
	"sync"

	"com.hanacaraka/domain/entities"
	"com.hanacaraka/domain/repositories"
)

// MemoryUserRepository implements UserRepository interface using in-memory storage
type MemoryUserRepository struct {
	users []*entities.User
	mutex sync.RWMutex
}

// NewMemoryUserRepository creates a new MemoryUserRepository instance
func NewMemoryUserRepository() repositories.UserRepository {
	repo := &MemoryUserRepository{
		users: make([]*entities.User, 0),
	}

	// Initialize with some sample data
	repo.users = append(repo.users, entities.NewUser("John Doe", "john@example.com"))
	repo.users = append(repo.users, entities.NewUser("Jane Smith", "jane@example.com"))
	repo.users = append(repo.users, entities.NewUser("Bob Johnson", "bob@example.com"))

	return repo
}

// GetAll returns all users
func (r *MemoryUserRepository) GetAll() ([]*entities.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	result := make([]*entities.User, len(r.users))
	copy(result, r.users)
	return result, nil
}

// GetByID returns a user by ID
func (r *MemoryUserRepository) GetByID(id string) (*entities.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	for _, user := range r.users {
		if user.ID == id {
			return user, nil
		}
	}
	return nil, errors.New("user not found")
}

// Create creates a new user
func (r *MemoryUserRepository) Create(user *entities.User) (*entities.User, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.users = append(r.users, user)
	return user, nil
}

// Update updates an existing user
func (r *MemoryUserRepository) Update(user *entities.User) (*entities.User, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	for i, existingUser := range r.users {
		if existingUser.ID == user.ID {
			r.users[i] = user
			return user, nil
		}
	}
	return nil, errors.New("user not found")
}

// Delete deletes a user by ID
func (r *MemoryUserRepository) Delete(id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	for i, user := range r.users {
		if user.ID == id {
			r.users = append(r.users[:i], r.users[i+1:]...)
			return nil
		}
	}
	return errors.New("user not found")
}
