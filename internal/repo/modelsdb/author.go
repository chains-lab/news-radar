package modelsdb

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Author struct {
	ID        uuid.UUID  `json:"id" bson:"id"`
	Name      string     `json:"name" bson:"name"`
	Status    string     `json:"status" bson:"status"`
	Desc      *string    `json:"desc" bson:"desc"`
	Avatar    *string    `json:"avatar,omitempty" bson:"avatar,omitempty"`
	Email     *string    `json:"email,omitempty" bson:"email,omitempty"`
	Telegram  *string    `json:"telegram,omitempty" bson:"telegram,omitempty"`
	Twitter   *string    `json:"twitter,omitempty" bson:"twitter,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	CreatedAt time.Time  `json:"created_at" bson:"created_at"`
}

// AuthorMongo is the model for the author in the MongoDB
type AuthorMongo struct {
	ID        uuid.UUID  `json:"id" bson:"id"`
	Name      string     `json:"name" bson:"name"`
	Desc      *string    `json:"desc" bson:"desc"`
	Avatar    *string    `json:"avatar,omitempty" bson:"avatar,omitempty"`
	Email     *string    `json:"email,omitempty" bson:"email,omitempty"`
	Telegram  *string    `json:"telegram,omitempty" bson:"telegram,omitempty"`
	Twitter   *string    `json:"twitter,omitempty" bson:"twitter,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	CreatedAt time.Time  `json:"created_at" bson:"created_at"`
}

// AuthorNeo is the model for the author in the Neo4j
type AuthorNeo struct {
	ID     uuid.UUID `json:"id"`
	Name   string    `json:"name"`
	Status string    `json:"status"`
}

func AuthorsCreateModel(mongo AuthorMongo, neo AuthorNeo) (Author, error) {
	if mongo.ID != neo.ID {
		return Author{}, fmt.Errorf("mongo and neo IDs do not match")
	}

	return Author{
		ID:        mongo.ID,
		Name:      mongo.Name,
		Status:    neo.Status,
		Desc:      mongo.Desc,
		Avatar:    mongo.Avatar,
		Email:     mongo.Email,
		Telegram:  mongo.Telegram,
		Twitter:   mongo.Twitter,
		UpdatedAt: mongo.UpdatedAt,
		CreatedAt: mongo.CreatedAt,
	}, nil
}
