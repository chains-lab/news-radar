package domain

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/recovery-flow/news-radar/internal/service/domain/models"
	"github.com/recovery-flow/news-radar/internal/service/infra/data/neo"
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

	AddAuthor(ctx context.Context, ID uuid.UUID, author uuid.UUID) error
	DeleteAuthor(ctx context.Context, ID uuid.UUID) error
	SetAuthors(ctx context.Context, ID uuid.UUID, author []uuid.UUID) error

	GetByID(ctx context.Context, ID uuid.UUID) (*models.Article, error)
}

func (d *domain) Create(
	ctx context.Context,
	title string,
	icon string,
	desc string,
	content []models.Section,
) (*models.Article, error) {
	articleID := uuid.New()

	err := d.Infra.Neo.Articles.Create(ctx, &neo.Article{
		ID:        articleID,
		CreatedAt: time.Now().UTC(),
		Tags:      nil,
		Themes:    nil,
	})
	if err != nil {
		return nil, err
	}

	res, err := d.Infra.Mongo.Articles.Insert(ctx, &models.Article{
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

func (d *domain) Update(
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

	res, err := d.Infra.Mongo.Articles.FilterID(ID).Update(ctx, fields)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (d *domain) Delete(ID uuid.UUID, ctx context.Context) error {
	err := d.Infra.Neo.Articles.Delete(ctx, ID)
	if err != nil {
		return err
	}

	err = d.Infra.Mongo.Articles.FilterID(ID).Delete(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (d *domain) GetByID(ctx context.Context, ID uuid.UUID) (*models.Article, error) {
	res, err := d.Infra.Mongo.Articles.FilterID(ID).Get(ctx)
	if err != nil {
		return nil, err
	}

	return res, nil
}

//Tags

func (d *domain) SetTags(ctx context.Context, ID uuid.UUID, tags []string) error {
	for _, tag := range tags {
		curTag, err := d.Infra.Neo.Tags.FindByName(ctx, tag)
		if err != nil {
			return err
		}
		if curTag == nil {
			return fmt.Errorf("tag %s not found", tag)
		}
	}

	err := d.Infra.Neo.Articles.SetHasTag(ctx, ID, tags)
	if err != nil {
		return err
	}

	return nil
}

func (d *domain) AddTag(ctx context.Context, ID uuid.UUID, tag string) error {
	curTag, err := d.Infra.Neo.Tags.FindByName(ctx, tag)
	if err != nil {
		return err
	}
	if curTag == nil {
		return fmt.Errorf("tag %s not found", tag)
	}

	err = d.Infra.Neo.Articles.CreateHasTagRelationship(ctx, ID, tag)
	if err != nil {
		return err
	}

	return nil
}

func (d *domain) DeleteTag(ctx context.Context, ID uuid.UUID, tag string) error {
	curTag, err := d.Infra.Neo.Tags.FindByName(ctx, tag)
	if err != nil {
		return err
	}
	if curTag == nil {
		return fmt.Errorf("tag %s not found", tag)
	}

	err = d.Infra.Neo.Articles.DeleteHasTagRelationship(ctx, ID, tag)
	if err != nil {
		return err
	}

	return nil
}

//Theme

func (d *domain) SetTheme(ctx context.Context, ID uuid.UUID, theme []string) error {
	for _, them := range theme {
		curTheme, err := d.Infra.Neo.Themes.FindByName(ctx, them)
		if err != nil {
			return err
		}
		if curTheme == nil {
			return fmt.Errorf("theme %s not found", theme)
		}
	}

	err := d.Infra.Neo.Articles.SetAbout(ctx, ID, theme)
	if err != nil {
		return err
	}

	return nil
}

func (d *domain) AddTheme(ctx context.Context, ID uuid.UUID, theme string) error {
	curTheme, err := d.Infra.Neo.Themes.FindByName(ctx, theme)
	if err != nil {
		return err
	}
	if curTheme == nil {
		return fmt.Errorf("theme %s not found", theme)
	}

	err = d.Infra.Neo.Articles.CreateAboutRelationship(ctx, ID, theme)
	if err != nil {
		return err
	}

	return nil
}

func (d *domain) DeleteTheme(ctx context.Context, ID uuid.UUID, theme string) error {
	curTheme, err := d.Infra.Neo.Themes.FindByName(ctx, theme)
	if err != nil {
		return err
	}
	if curTheme == nil {
		return fmt.Errorf("theme %s not found", theme)
	}

	err = d.Infra.Neo.Articles.DeleteAboutRelationship(ctx, ID, theme)
	if err != nil {
		return err
	}

	return nil
}

//Author

func (d *domain) SetAuthors(ctx context.Context, ID uuid.UUID, authors []string) error {
	var IDs []uuid.UUID
	for _, author := range authors {
		ID, err := uuid.Parse(author)
		if err != nil {
			return err
		}
		curAuthor, err := d.Infra.Neo.Authors.GetByID(ctx, ID)
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

	err := d.Infra.Neo.Articles.SetAuthors(ctx, ID, IDs)
	if err != nil {
		return err
	}

	return nil
}

func (d *domain) AddAuthor(ctx context.Context, ID uuid.UUID, author uuid.UUID) error {
	curAuthor, err := d.Infra.Neo.Authors.GetByID(ctx, author)
	if err != nil {
		return err
	}
	if curAuthor == nil {
		return fmt.Errorf("author %s not found", author)
	}

	err = d.Infra.Neo.Articles.CreateAuthorshipRelationship(ctx, ID, author)
	if err != nil {
		return err
	}

	return nil
}

func (d *domain) DeleteAuthor(ctx context.Context, ID uuid.UUID) error {
	curAuthor, err := d.Infra.Neo.Authors.GetByID(ctx, ID)
	if err != nil {
		return err
	}
	if curAuthor == nil {
		return fmt.Errorf("author %s not found", ID)
	}

	err = d.Infra.Neo.Articles.DeleteAuthorshipRelationship(ctx, ID, ID)
	if err != nil {
		return err
	}

	return nil
}
