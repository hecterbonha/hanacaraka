package entities

import "github.com/google/uuid"

type Player struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Status Status `json:"status"`
}

type Status struct {
	STR int `json:"STR"`
}

func NewPlayer(name, email string) *Player {
	return &Player{
		ID:    uuid.New().String(),
		Name:  name,
	}
}

func NewPlayerWithID(id, name, email string) *Player {
	return &Player{
		ID:    id,
		Name:  name,
	}
}

func (u *Player) IsValid() bool {
	return u.Name != ""
}
