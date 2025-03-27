package entities

import (
	"context"

	"github.com/google/uuid"
	"github.com/recovery-flow/news-radar/internal/config"
	"github.com/recovery-flow/news-radar/internal/data"
)

type reactionRepo interface {
	CreateLike(ctx context.Context, userID, articleID uuid.UUID) error
	RemoveLike(ctx context.Context, userID, articleID uuid.UUID) error

	CreateDislike(ctx context.Context, userID, articleID uuid.UUID) error
	RemoveDislike(ctx context.Context, userID, articleID uuid.UUID) error

	CreateRepost(ctx context.Context, userID, articleID uuid.UUID) error
}

type Reaction struct {
	data reactionRepo
}

func NewReaction(cfg config.Config) (*Reaction, error) {
	repo, err := data.NewReactions(cfg)
	if err != nil {
		return nil, err
	}

	return &Reaction{
		data: repo,
	}, nil
}

func (r *Reaction) MakeLike(ctx context.Context, userID, articleID uuid.UUID) error {
	return r.data.CreateLike(ctx, userID, articleID)
}

func (r *Reaction) RemoveLike(ctx context.Context, userID, articleID uuid.UUID) error {
	return r.data.RemoveLike(ctx, userID, articleID)
}

func (r *Reaction) MakeDislike(ctx context.Context, userID, articleID uuid.UUID) error {
	return r.data.CreateDislike(ctx, userID, articleID)
}

func (r *Reaction) RemoveDislike(ctx context.Context, userID, articleID uuid.UUID) error {
	return r.data.RemoveDislike(ctx, userID, articleID)
}

func (r *Reaction) MakeRepost(ctx context.Context, userID, articleID uuid.UUID) error {
	return r.data.CreateRepost(ctx, userID, articleID)
}
