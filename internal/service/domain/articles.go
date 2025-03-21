package domain

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/recovery-flow/news-radar/internal/service/infra"
	"github.com/recovery-flow/news-radar/internal/service/infra/data/neo"
	"github.com/recovery-flow/news-radar/internal/service/models"
	"github.com/sirupsen/logrus"
)

type Articles interface {
	Create(
		ctx context.Context,
		title string,
		icon string,
		desc string,
		content []models.Section,
	) (*models.Article, error)
	Update(
		ctx context.Context,
		ID uuid.UUID,
		title *string,
		icon *string,
		desc *string,
		status *string,
		content []models.Section,
	) (*models.Article, error)
	Delete(ID uuid.UUID, ctx context.Context) error

	SetTags(ctx context.Context, ID uuid.UUID, tags []string) error
	AddTag(ctx context.Context, ID uuid.UUID, tag string) error
	DeleteTag(ctx context.Context, ID uuid.UUID, tag string) error

	SetTheme(ctx context.Context, ID uuid.UUID, theme []string) error
	AddTheme(ctx context.Context, ID uuid.UUID, theme string) error
	DeleteTheme(ctx context.Context, ID uuid.UUID, theme string) error

	AddAuthor(ctx context.Context, ID uuid.UUID, author string) error
	DeleteAuthor(ctx context.Context, ID uuid.UUID, author string) error
	SetAuthors(ctx context.Context, ID uuid.UUID, authors []string) error

	GetByID(ctx context.Context, ID uuid.UUID) (*models.Article, error)
}

type articles struct {
	Infra *infra.Infra
	log   *logrus.Logger
}

func (a *articles) Create(
	ctx context.Context,
	title string,
	icon string,
	desc string,
	content []models.Section,
) (*models.Article, error) {
	articleID := uuid.New()

	err := a.Infra.Neo.Articles.Create(ctx, &neo.Article{
		ID:        articleID,
		CreatedAt: time.Now().UTC(),
		Tags:      nil,
		Themes:    nil,
	})
	if err != nil {
		return nil, err
	}

	res, err := a.Infra.Mongo.Articles.Insert(ctx, &models.Article{
		ID:        articleID,
		Title:     title,
		Icon:      icon,
		Desc:      desc,
		Content:   content,
		Likes:     0,
		Reposts:   0,
		Status:    models.ArticleStatusInactive,
		CreatedAt: time.Now().UTC(),
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (a *articles) Update(
	ctx context.Context,
	ID uuid.UUID,
	title *string,
	icon *string,
	desc *string,
	status *string,
	content []models.Section,
) (*models.Article, error) {
	fields := make(map[string]any)
	if title != nil {
		fields["title"] = title
	}
	if icon != nil {
		fields["icon"] = icon
	}
	if desc != nil {
		fields["desc"] = desc
	}
	if content != nil {
		fields["content"] = content
	}
	if status != nil {
		st, err := models.ParseArticleStatus(*status)
		if err != nil {
			return nil, err
		}
		fields["status"] = st
	}

	res, err := a.Infra.Mongo.Articles.FilterID(ID).Update(ctx, fields)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (a *articles) Delete(ID uuid.UUID, ctx context.Context) error {
	err := a.Infra.Neo.Articles.Delete(ctx, ID)
	if err != nil {
		return err
	}

	err = a.Infra.Mongo.Articles.FilterID(ID).Delete(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (a *articles) GetByID(ctx context.Context, ID uuid.UUID) (*models.Article, error) {
	res, err := a.Infra.Mongo.Articles.FilterID(ID).Get(ctx)
	if err != nil {
		return nil, err
	}

	return res, nil
}

//Tags

func (a *articles) SetTags(ctx context.Context, ID uuid.UUID, tags []string) error {
	for _, tag := range tags {
		curTag, err := a.Infra.Neo.Tags.FindByName(ctx, tag)
		if err != nil {
			return err
		}
		if curTag == nil {
			return fmt.Errorf("tag %s not found", tag)
		}
	}

	err := a.Infra.Neo.Articles.SetHasTag(ctx, ID, tags)
	if err != nil {
		return err
	}

	return nil
}

func (a *articles) AddTag(ctx context.Context, ID uuid.UUID, tag string) error {
	curTag, err := a.Infra.Neo.Tags.FindByName(ctx, tag)
	if err != nil {
		return err
	}
	if curTag == nil {
		return fmt.Errorf("tag %s not found", tag)
	}

	err = a.Infra.Neo.Articles.CreateHasTagRelationship(ctx, ID, tag)
	if err != nil {
		return err
	}

	return nil
}

func (a *articles) DeleteTag(ctx context.Context, ID uuid.UUID, tag string) error {
	curTag, err := a.Infra.Neo.Tags.FindByName(ctx, tag)
	if err != nil {
		return err
	}
	if curTag == nil {
		return fmt.Errorf("tag %s not found", tag)
	}

	err = a.Infra.Neo.Articles.DeleteHasTagRelationship(ctx, ID, tag)
	if err != nil {
		return err
	}

	return nil
}

//Theme

func (a *articles) SetTheme(ctx context.Context, ID uuid.UUID, theme []string) error {
	for _, them := range theme {
		curTheme, err := a.Infra.Neo.Themes.FindByName(ctx, them)
		if err != nil {
			return err
		}
		if curTheme == nil {
			return fmt.Errorf("theme %s not found", theme)
		}
	}

	err := a.Infra.Neo.Articles.SetAbout(ctx, ID, theme)
	if err != nil {
		return err
	}

	return nil
}

func (a *articles) AddTheme(ctx context.Context, ID uuid.UUID, theme string) error {
	curTheme, err := a.Infra.Neo.Themes.FindByName(ctx, theme)
	if err != nil {
		return err
	}
	if curTheme == nil {
		return fmt.Errorf("theme %s not found", theme)
	}

	err = a.Infra.Neo.Articles.CreateAboutRelationship(ctx, ID, theme)
	if err != nil {
		return err
	}

	return nil
}

func (a *articles) DeleteTheme(ctx context.Context, ID uuid.UUID, theme string) error {
	curTheme, err := a.Infra.Neo.Themes.FindByName(ctx, theme)
	if err != nil {
		return err
	}
	if curTheme == nil {
		return fmt.Errorf("theme %s not found", theme)
	}

	err = a.Infra.Neo.Articles.DeleteAboutRelationship(ctx, ID, theme)
	if err != nil {
		return err
	}

	return nil
}

//Author

func (a *articles) SetAuthors(ctx context.Context, ID uuid.UUID, authors []string) error {
	var IDs []uuid.UUID
	for _, author := range authors {
		ID, err := uuid.Parse(author)
		if err != nil {
			return err
		}
		curAuthor, err := a.Infra.Neo.Authors.GetByID(ctx, ID)
		if err != nil {
			return err
		}
		if curAuthor == nil {
			return fmt.Errorf("author %s not found", ID)
		}
		IDs = append(IDs, ID)
	}
	if IDs == nil {
		return nil
	}

	err := a.Infra.Neo.Articles.SetAuthors(ctx, ID, IDs)
	if err != nil {
		return err
	}

	return nil
}

func (a *articles) AddAuthor(ctx context.Context, ID uuid.UUID, author string) error {
	authorID, err := uuid.Parse(author)
	if err != nil {
		return err
	}
	curAuthor, err := a.Infra.Neo.Authors.GetByID(ctx, authorID)
	if err != nil {
		return err
	}
	if curAuthor == nil {
		return fmt.Errorf("author %s not found", authorID)
	}

	err = a.Infra.Neo.Articles.CreateAuthorshipRelationship(ctx, ID, authorID)
	if err != nil {
		return err
	}

	return nil
}

func (a *articles) DeleteAuthor(ctx context.Context, ID uuid.UUID, author string) error {
	authorID, err := uuid.Parse(author)
	if err != nil {
		return err
	}
	curAuthor, err := a.Infra.Neo.Authors.GetByID(ctx, ID)
	if err != nil {
		return err
	}
	if curAuthor == nil {
		return fmt.Errorf("author %s not found", authorID)
	}

	err = a.Infra.Neo.Articles.DeleteAuthorshipRelationship(ctx, ID, authorID)
	if err != nil {
		return err
	}

	return nil
}
