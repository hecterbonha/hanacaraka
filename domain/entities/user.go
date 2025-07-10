package entities

import "github.com/google/uuid"

type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func NewUser(name, email string) *User {
	return &User{
		ID:    uuid.New().String(),
		Name:  name,
		Email: email,
	}
}

func NewUserWithID(id, name, email string) *User {
	return &User{
		ID:    id,
		Name:  name,
		Email: email,
	}
}

func (u *User) IsValid() bool {
	return u.Name != "" && u.Email != ""
}

func (u *User) UpdateName(name string) {
	u.Name = name
}

func (u *User) UpdateEmail(email string) {
	u.Email = email
}
