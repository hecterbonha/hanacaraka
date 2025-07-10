package repositories

import "com.hanacaraka/domain/entities"

type UserRepository interface {
	GetAll() ([]*entities.User, error)
	GetByID(id string) (*entities.User, error)
	Create(user *entities.User) (*entities.User, error)
	Update(user *entities.User) (*entities.User, error)
	Delete(id string) error
}
