package services

import (
	"errors"

	"com.hanacaraka/domain/entities"
	"com.hanacaraka/domain/repositories"
)

type UserServiceInterface interface {
	GetAllUsers() ([]*entities.User, error)
	GetUserByID(id string) (*entities.User, error)
	CreateUser(name, email string) (*entities.User, error)
	UpdateUser(id string, name, email string) (*entities.User, error)
	DeleteUser(id string) error
}

type UserService struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) GetAllUsers() ([]*entities.User, error) {
	return s.userRepo.GetAll()
}

func (s *UserService) GetUserByID(id string) (*entities.User, error) {
	if id == "" {
		return nil, errors.New("invalid user ID")
	}
	return s.userRepo.GetByID(id)
}

func (s *UserService) CreateUser(name, email string) (*entities.User, error) {
	user := entities.NewUser(name, email)

	if !user.IsValid() {
		return nil, errors.New("invalid user data: name and email are required")
	}

	return s.userRepo.Create(user)
}

func (s *UserService) UpdateUser(id string, name, email string) (*entities.User, error) {
	if id == "" {
		return nil, errors.New("invalid user ID")
	}

	if name == "" || email == "" {
		return nil, errors.New("invalid user data: name and email are required")
	}

	existingUser, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if existingUser == nil {
		return nil, errors.New("user not found")
	}

	existingUser.UpdateName(name)
	existingUser.UpdateEmail(email)

	if !existingUser.IsValid() {
		return nil, errors.New("invalid user data: name and email are required")
	}

	return s.userRepo.Update(existingUser)
}

func (s *UserService) DeleteUser(id string) error {
	if id == "" {
		return errors.New("invalid user ID")
	}

	existingUser, err := s.userRepo.GetByID(id)
	if err != nil {
		return err
	}

	if existingUser == nil {
		return errors.New("user not found")
	}

	return s.userRepo.Delete(id)
}
