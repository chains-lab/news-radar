package app

import (
	"github.com/google/uuid"
	"github.com/hs-zavet/news-radar/internal/config"
	"github.com/hs-zavet/news-radar/internal/repo"
	"github.com/sirupsen/logrus"
)

type App struct {
	articles articlesRepo
	authors  authorsRepo
	tags     tagsRepo
	log      *logrus.Entry
}

func NewApp(cfg config.Config, log *logrus.Logger) (App, error) {
	articles, err := newArticles(cfg)
	if err != nil {
		return App{}, err
	}
	authors, err := newAuthors(cfg)
	if err != nil {
		return App{}, err
	}
	tags, err := newTags(cfg)
	if err != nil {
		return App{}, err
	}

	return App{
		articles: articles,
		authors:  authors,
		tags:     tags,
		log:      log.WithField("component", "app"),
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
