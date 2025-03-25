package entities

import (
	"github.com/recovery-flow/news-radar/internal/config"
	"github.com/recovery-flow/news-radar/internal/events/writer"
	"github.com/recovery-flow/news-radar/internal/service/repo"
)

type User interface {
}

type user struct {
	data   repo.Users
	writer writer.Reaction
}

func NewUser(cfg config.Config) (User, error) {
	data, err := repo.NewUsers(cfg)
	if err != nil {
		return nil, err
	}

	reaction := writer.NewReactions(cfg)

	return &user{
		data:   data,
		writer: reaction,
	}, nil
}
