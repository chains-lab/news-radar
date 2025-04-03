package repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/hs-zavet/news-radar/internal/config"
	"github.com/hs-zavet/news-radar/internal/repo/neodb"
)

type reaction interface {
	Create(ctx context.Context, userID uuid.UUID, articleID uuid.UUID) error
	Delete(ctx context.Context, userID uuid.UUID, articleID uuid.UUID) error

	GetForUser(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error)
	GetForArticle(ctx context.Context, articleID uuid.UUID) ([]uuid.UUID, error)
	GetForUserAndArticle(ctx context.Context, articleID, userID uuid.UUID) (bool, error)
}

type Reaction struct {
	likes   reaction
	reposts reaction
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

	return &Reaction{
		likes:   likes,
		reposts: reposts,
	}, nil
}

func (u *Reaction) CreateLike(userID, articleID uuid.UUID) error {
	ctxSync, cancel := context.WithTimeout(context.Background(), dataCtxTimeAisle)
	defer cancel()

	return u.likes.Create(ctxSync, userID, articleID)
}

func (u *Reaction) GetLikesForUserAndArticle(userID, articleID uuid.UUID) (bool, error) {
	ctxSync, cancel := context.WithTimeout(context.Background(), dataCtxTimeAisle)
	defer cancel()

	return u.likes.GetForUserAndArticle(ctxSync, articleID, userID)
}

func (u *Reaction) RemoveLike(userID, articleID uuid.UUID) error {
	ctxSync, cancel := context.WithTimeout(context.Background(), dataCtxTimeAisle)
	defer cancel()

	return u.likes.Delete(ctxSync, userID, articleID)
}

func (u *Reaction) CreateRepost(userID, articleID uuid.UUID) error {
	ctxSync, cancel := context.WithTimeout(context.Background(), dataCtxTimeAisle)
	defer cancel()

	return u.reposts.Create(ctxSync, userID, articleID)
}

func (u *Reaction) GetRepostsForUserAndArticle(userID, articleID uuid.UUID) (bool, error) {
	ctxSync, cancel := context.WithTimeout(context.Background(), dataCtxTimeAisle)
	defer cancel()

	return u.reposts.GetForUserAndArticle(ctxSync, articleID, userID)
}
