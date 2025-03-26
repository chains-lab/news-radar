package repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/recovery-flow/news-radar/internal/config"
	"github.com/recovery-flow/news-radar/internal/data/neodb"
)

type Users interface {
	Create(ctx context.Context, userID uuid.UUID) error
	Get(ctx context.Context, userID uuid.UUID) (*neodb.User, error)

	AddLike(ctx context.Context, userID, articleID uuid.UUID) error
	RemoveLike(ctx context.Context, userID, articleID uuid.UUID) error

	AddDislike(ctx context.Context, userID, articleID uuid.UUID) error
	RemoveDislike(ctx context.Context, userID, articleID uuid.UUID) error

	AddRepost(ctx context.Context, userID, articleID uuid.UUID) error
}

type users struct {
	neo      neodb.Users
	likes    neodb.Likes
	reposts  neodb.Reposts
	dislikes neodb.Dislikes
}

func NewUsers(cfg config.Config) (Users, error) {
	neo, err := neodb.NewUsers(cfg.Database.Neo4j.URI, cfg.Database.Neo4j.Username, cfg.Database.Neo4j.Password)
	if err != nil {
		return nil, err
	}

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

	return &users{
		neo:      neo,
		likes:    likes,
		reposts:  reposts,
		dislikes: dislikes,
	}, nil
}

func (u *users) Create(ctx context.Context, userID uuid.UUID) error {
	return u.neo.Create(ctx, neodb.User{
		ID: userID,
	})
}

func (u *users) Get(ctx context.Context, userID uuid.UUID) (*neodb.User, error) {
	return u.neo.Get(ctx, userID)
}

func (u *users) AddLike(ctx context.Context, userID, articleID uuid.UUID) error {
	return u.likes.Create(ctx, userID, articleID)
}

func (u *users) RemoveLike(ctx context.Context, userID, articleID uuid.UUID) error {
	return u.likes.Delete(ctx, userID, articleID)
}

func (u *users) AddDislike(ctx context.Context, userID, articleID uuid.UUID) error {
	return u.dislikes.Create(ctx, userID, articleID)
}

func (u *users) RemoveDislike(ctx context.Context, userID, articleID uuid.UUID) error {
	return u.dislikes.Delete(ctx, userID, articleID)
}

func (u *users) AddRepost(ctx context.Context, userID, articleID uuid.UUID) error {
	return u.reposts.Create(ctx, userID, articleID)
}
