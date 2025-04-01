package modelsdb

import (
	"github.com/google/uuid"
)

type User struct {
	ID uuid.UUID `json:"id"`
}

// UserNeo is the model for the user in the Neo4j
type UserNeo struct {
	ID uuid.UUID `json:"id"`
}

func NewUser(neo UserNeo) User {
	return User{
		ID: neo.ID,
	}
}
