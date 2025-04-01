package models

import (
	"github.com/google/uuid"
	"github.com/hs-zavet/news-radar/internal/data/neodb"
)

type User struct {
	ID uuid.UUID `json:"id"`
}

func NewUser(neo neodb.UserModels) User {
	return User{
		ID: neo.ID,
	}
}
