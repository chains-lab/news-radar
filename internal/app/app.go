package app

import (
	"github.com/google/uuid"
	"github.com/hs-zavet/news-radar/internal/config"
	"github.com/hs-zavet/news-radar/internal/repo"
)

type App struct {
	articles  articlesRepo
	authors   authorsRepo
	reactions reactionsRepo
	tags      tagsRepo
	users     usersRepo
}

func NewApp(cfg config.Config) (App, error) {
	articles, err := newArticles(cfg)
	if err != nil {
		return App{}, err
	}
	authors, err := newAuthors(cfg)
	if err != nil {
		return App{}, err
	}
	reactions, err := newReactions(cfg)
	if err != nil {
		return App{}, err
	}
	tags, err := newTags(cfg)
	if err != nil {
		return App{}, err
	}
	users, err := newUsers(cfg)
	if err != nil {
		return App{}, err
	}

	return App{
		articles:  articles,
		authors:   authors,
		reactions: reactions,
		tags:      tags,
		users:     users,
	}, nil
}

type articlesRepo interface {
	Create(input repo.ArticleCreateInput) error
	Update(ID uuid.UUID, input repo.ArticleUpdateInput) error
	Delete(ID uuid.UUID) error

	SetTags(ID uuid.UUID, tags []string) error
	AddTag(ID uuid.UUID, tag string) error
	DeleteTag(ID uuid.UUID, tag string) error
	GetTags(ID uuid.UUID) ([]string, error)

	AddAuthor(ID uuid.UUID, author uuid.UUID) error
	DeleteAuthor(ID uuid.UUID, author uuid.UUID) error
	SetAuthors(ID uuid.UUID, authors []uuid.UUID) error
	GetAuthors(ID uuid.UUID) ([]uuid.UUID, error)

	GetByID(ID uuid.UUID) (repo.ArticleModel, error)
}

func newArticles(cfg config.Config) (articlesRepo, error) {
	data, err := repo.NewArticles(cfg)
	if err != nil {
		return nil, err
	}
	return data, nil
}

type authorsRepo interface {
	Create(author repo.AuthorCreateInput) error
	Update(ID uuid.UUID, input repo.AuthorUpdateInput) error
	Delete(ID uuid.UUID) error

	GetByID(ID uuid.UUID) (repo.AuthorModel, error)
}

func newAuthors(cfg config.Config) (authorsRepo, error) {
	data, err := repo.NewAuthors(cfg)
	if err != nil {
		return nil, err
	}
	return data, nil
}

type reactionsRepo interface {
	CreateLike(userID, articleID uuid.UUID) error
	RemoveLike(userID, articleID uuid.UUID) error
	GetRepostsForUserAndArticle(userID, articleID uuid.UUID) (bool, error)
	GetLikesForUserAndArticle(userID, articleID uuid.UUID) (bool, error)

	CreateRepost(userID, articleID uuid.UUID) error
}

func newReactions(cfg config.Config) (reactionsRepo, error) {
	data, err := repo.NewReactions(cfg)
	if err != nil {
		return nil, err
	}
	return data, nil
}

type tagsRepo interface {
	Create(input repo.TagCreateInput) error
	Delete(name string) error
	Update(name string, input repo.TagUpdateInput) error
	Get(name string) (repo.TagModel, error)
}

func newTags(cfg config.Config) (tagsRepo, error) {
	data, err := repo.NewTags(cfg)
	if err != nil {
		return nil, err
	}
	return data, nil
}

type usersRepo interface {
	Create(userID repo.UserCreateInput) error
	Get(userID uuid.UUID) (repo.UserModel, error)
}

func newUsers(cfg config.Config) (usersRepo, error) {
	data, err := repo.NewUsers(cfg)
	if err != nil {
		return nil, err
	}
	return data, nil
}
