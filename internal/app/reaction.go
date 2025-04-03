package domain

import (
	"context"

	"github.com/google/uuid"
	"github.com/hs-zavet/news-radar/internal/config"
	"github.com/hs-zavet/news-radar/internal/repo"
)

type reactionRepo interface {
	CreateLike(userID, articleID uuid.UUID) error
	RemoveLike(userID, articleID uuid.UUID) error

	CreateDislike(userID, articleID uuid.UUID) error
	RemoveDislike(userID, articleID uuid.UUID) error

	CreateRepost(userID, articleID uuid.UUID) error
}

type Reaction struct {
	data reactionRepo
}

func NewReaction(cfg config.Config) (*Reaction, error) {
	repo, err := repo.NewReactions(cfg)
	if err != nil {
		return nil, err
	}

	return &Reaction{
		data: repo,
	}, nil
}

func (r *Reaction) AddLike(ctx context.Context, userID, articleID uuid.UUID) error {
	return r.data.CreateLike(userID, articleID)
}

func (r *Reaction) RemoveLike(ctx context.Context, userID, articleID uuid.UUID) error {
	return r.data.RemoveLike(userID, articleID)
}

func (r *Reaction) AddDislike(ctx context.Context, userID, articleID uuid.UUID) error {
	return r.data.CreateDislike(userID, articleID)
}

func (r *Reaction) RemoveDislike(ctx context.Context, userID, articleID uuid.UUID) error {
	return r.data.RemoveDislike(userID, articleID)
}

func (r *Reaction) AddRepost(ctx context.Context, userID, articleID uuid.UUID) error {
	return r.data.CreateRepost(userID, articleID)
}
