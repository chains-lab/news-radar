package repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/recovery-flow/news-radar/internal/config"
	"github.com/recovery-flow/news-radar/internal/data/neodb"
)

type UsersNeo interface {
	Create(ctx context.Context, user neodb.UserModels) error
	Get(ctx context.Context, id uuid.UUID) (*neodb.UserModels, error)
}

type Likes interface {
	Create(ctx context.Context, userID uuid.UUID, articleID uuid.UUID) error
	Delete(ctx context.Context, userID uuid.UUID, articleID uuid.UUID) error

	GetForUser(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error)
	GetForArticle(ctx context.Context, articleID uuid.UUID) ([]uuid.UUID, error)
}

type Reposts interface {
	Create(ctx context.Context, userID uuid.UUID, articleID uuid.UUID) error

	GetForUser(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error)
	GetForArticle(ctx context.Context, articleID uuid.UUID) ([]uuid.UUID, error)
}

type Dislikes interface {
	Create(ctx context.Context, userID uuid.UUID, articleID uuid.UUID) error
	Delete(ctx context.Context, userID uuid.UUID, articleID uuid.UUID) error

	GetForUser(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error)
	GetForArticle(ctx context.Context, articleID uuid.UUID) ([]uuid.UUID, error)
}

type Users struct {
	neo      UsersNeo
	likes    Likes
	reposts  Reposts
	dislikes Dislikes
}

func NewUsers(cfg config.Config) (*Users, error) {
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

	return &Users{
		neo:      neo,
		likes:    likes,
		reposts:  reposts,
		dislikes: dislikes,
	}, nil
}

func (u *Users) Create(ctx context.Context, userID uuid.UUID) error {
	return u.neo.Create(ctx, neodb.UserModels{
		ID: userID,
	})
}

func (u *Users) Get(ctx context.Context, userID uuid.UUID) (*neodb.UserModels, error) {
	return u.neo.Get(ctx, userID)
}

func (u *Users) AddLike(ctx context.Context, userID, articleID uuid.UUID) error {
	return u.likes.Create(ctx, userID, articleID)
}

func (u *Users) RemoveLike(ctx context.Context, userID, articleID uuid.UUID) error {
	return u.likes.Delete(ctx, userID, articleID)
}

func (u *Users) AddDislike(ctx context.Context, userID, articleID uuid.UUID) error {
	return u.dislikes.Create(ctx, userID, articleID)
}

func (u *Users) RemoveDislike(ctx context.Context, userID, articleID uuid.UUID) error {
	return u.dislikes.Delete(ctx, userID, articleID)
}

func (u *Users) AddRepost(ctx context.Context, userID, articleID uuid.UUID) error {
	return u.reposts.Create(ctx, userID, articleID)
}
