package entities

import (
	"github.com/recovery-flow/news-radar/internal/config"
	"github.com/recovery-flow/news-radar/internal/service/data"
	"github.com/recovery-flow/news-radar/internal/service/events/writer"
)

type User interface {
}

type user struct {
	data   data.Users
	writer writer.Reaction
}

func NewUser(cfg config.Config) (User, error) {
	repo, err := data.NewUsers(cfg)
	if err != nil {
		return nil, err
	}

	reaction := writer.NewReactions(cfg)

	return &user{
		data:   repo,
		writer: reaction,
	}, nil
}
