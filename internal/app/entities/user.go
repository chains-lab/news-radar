package entities

import (
	"github.com/recovery-flow/news-radar/internal/config"
	"github.com/recovery-flow/news-radar/internal/events/writer"
	"github.com/recovery-flow/news-radar/internal/repo"
)

type User interface {
}

type user struct {
	data   repo.Users
	writer writer.Reaction
}

func NewUser(cfg config.Config) (User, error) {
	repo, err := repo.NewUsers(cfg)
	if err != nil {
		return nil, err
	}

	reaction := writer.NewReactions(cfg)

	return &user{
		data:   repo,
		writer: reaction,
	}, nil
}
