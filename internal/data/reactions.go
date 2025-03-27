package data

import (
	"context"

	"github.com/google/uuid"
	"github.com/recovery-flow/news-radar/internal/config"
	"github.com/recovery-flow/news-radar/internal/data/neodb"
)

type reactionRepo interface {
	Create(ctx context.Context, userID uuid.UUID, articleID uuid.UUID) error
	Delete(ctx context.Context, userID uuid.UUID, articleID uuid.UUID) error

	GetForUser(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error)
	GetForArticle(ctx context.Context, articleID uuid.UUID) ([]uuid.UUID, error)
}

type Reaction struct {
	likes    reactionRepo
	reposts  reactionRepo
	dislikes reactionRepo
}

func NewReactions(cfg config.Config) (*Reaction, error) {
	likes, err := neodb.NewLikes(cfg.Database.Neo4j.URI, cfg.Database.Neo4j.Username, cfg.Database.Neo4j.Password)
	if err != nil {
		return nil, err
	}

	reposts, err := neodb.NewReposts(cfg.Database.Neo4j.URI, cfg.Database.Neo4j.Username, cfg.Database.Neo4j.Password)
	if err != nil {
		return nil, err
	}

	dislikes, err := neodb.NewDislikes(cfg.Database.Neo4j.URI, cfg.Database.Neo4j.Username, cfg.Database.Neo4j.Password)
	if err != nil {
		return nil, err
	}

	return &Reaction{
		likes:    likes,
		reposts:  reposts,
		dislikes: dislikes,
	}, nil
}

func (u *Reaction) CreateLike(ctx context.Context, userID, articleID uuid.UUID) error {
	return u.likes.Create(ctx, userID, articleID)
}

func (u *Reaction) RemoveLike(ctx context.Context, userID, articleID uuid.UUID) error {
	return u.likes.Delete(ctx, userID, articleID)
}

func (u *Reaction) CreateDislike(ctx context.Context, userID, articleID uuid.UUID) error {
	return u.dislikes.Create(ctx, userID, articleID)
}

func (u *Reaction) RemoveDislike(ctx context.Context, userID, articleID uuid.UUID) error {
	return u.dislikes.Delete(ctx, userID, articleID)
}

func (u *Reaction) CreateRepost(ctx context.Context, userID, articleID uuid.UUID) error {
	return u.reposts.Create(ctx, userID, articleID)
}
