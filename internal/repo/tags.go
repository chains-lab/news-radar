package repo

import (
	"context"
	"time"

	"github.com/hs-zavet/news-radar/internal/config"
	"github.com/hs-zavet/news-radar/internal/enums"
	"github.com/hs-zavet/news-radar/internal/repo/neodb"
)

type TagModel struct {
	Name      string          `json:"name"`
	Status    enums.TagStatus `json:"status"`
	Type      enums.TagType   `json:"type"`
	Color     string          `json:"color"`
	Icon      string          `json:"icon"`
	UpdatedAt *time.Time      `json:"updated_at,omitempty"`
	CreatedAt time.Time       `json:"created_at"`
}

type tagsNeo interface {
	Create(ctx context.Context, input neodb.TagCreateInput) error
	Delete(ctx context.Context, name string) error

	UpdateStatus(ctx context.Context, name string, status enums.TagStatus, updatedAt time.Time) (neodb.TagModel, error)
	UpdateName(ctx context.Context, name string, newName string, updatedAt time.Time) (neodb.TagModel, error)
	UpdateType(ctx context.Context, name string, newType enums.TagType, updatedAt time.Time) (neodb.TagModel, error)
	UpdateColor(ctx context.Context, name string, color string, updatedAt time.Time) (neodb.TagModel, error)
	UpdateIcon(ctx context.Context, name string, icon string, updatedAt time.Time) (neodb.TagModel, error)

	Get(ctx context.Context, name string) (neodb.TagModel, error)
	GetAll(ctx context.Context) ([]neodb.TagModel, error)
}

type Tags struct {
	neo tagsNeo
}

func NewTags(cfg config.Config) (*Tags, error) {
	neo, err := neodb.NewTags(cfg.Database.Neo4j.URI, cfg.Database.Neo4j.Username, cfg.Database.Neo4j.Password)
	if err != nil {
		return nil, err
	}
	return &Tags{
		neo: neo,
	}, nil
}

type TagCreateInput struct {
	Name      string          `json:"name"`
	Status    enums.TagStatus `json:"status"`
	Type      enums.TagType   `json:"type"`
	Color     string          `json:"color"`
	Icon      string          `json:"icon"`
	CreatedAt time.Time       `json:"created_at"`
}

func (t *Tags) Create(input TagCreateInput) error {
	ctxSync, cancel := context.WithTimeout(context.Background(), dataCtxTimeAisle)
	defer cancel()

	err := t.neo.Create(ctxSync, neodb.TagCreateInput{
		Name:      input.Name,
		Status:    input.Status,
		Type:      input.Type,
		CreatedAt: input.CreatedAt,
	})
	if err != nil {
		return err
	}

	return nil
}

func (t *Tags) Delete(name string) error {
	ctxSync, cancel := context.WithTimeout(context.Background(), dataCtxTimeAisle)
	defer cancel()

	err := t.neo.Delete(ctxSync, name)
	if err != nil {
		return err
	}

	return nil
}

func (t *Tags) UpdateStatus(name string, status enums.TagStatus) (TagModel, error) {
	ctxSync, cancel := context.WithTimeout(context.Background(), dataCtxTimeAisle)
	defer cancel()

	res, err := t.neo.UpdateStatus(ctxSync, name, status)
	if err != nil {
		return TagModel{}, err
	}

	return TagModel{
		Name:      res.Name,
		Status:    res.Status,
		Type:      res.Type,
		Color:     res.Color,
		Icon:      res.Icon,
		CreatedAt: res.CreatedAt,
	}, nil
}

func (t *Tags) UpdateType(name string, tagType enums.TagType) (TagModel, error) {
	ctxSync, cancel := context.WithTimeout(context.Background(), dataCtxTimeAisle)
	defer cancel()

	res, err := t.neo.UpdateType(ctxSync, name, tagType)
	if err != nil {
		return TagModel{}, err
	}

	return TagModel{
		Name:      res.Name,
		Status:    res.Status,
		Type:      res.Type,
		Color:     res.Color,
		Icon:      res.Icon,
		CreatedAt: res.CreatedAt,
	}, nil
}

func (t *Tags) UpdateColor(name string, color string) (TagModel, error) {
	ctxSync, cancel := context.WithTimeout(context.Background(), dataCtxTimeAisle)
	defer cancel()

	res, err := t.neo.UpdateColor(ctxSync, name, color)
	if err != nil {
		return TagModel{}, err
	}

	return TagModel{
		Name:      res.Name,
		Status:    res.Status,
		Type:      res.Type,
		Color:     res.Color,
		Icon:      res.Icon,
		CreatedAt: res.CreatedAt,
	}, nil
}

func (t *Tags) UpdateIcon(name string, icon string) (TagModel, error) {
	ctxSync, cancel := context.WithTimeout(context.Background(), dataCtxTimeAisle)
	defer cancel()

	res, err := t.neo.UpdateIcon(ctxSync, name, icon)
	if err != nil {
		return TagModel{}, err
	}

	return TagModel{
		Name:      res.Name,
		Status:    res.Status,
		Type:      res.Type,
		Color:     res.Color,
		Icon:      res.Icon,
		CreatedAt: res.CreatedAt,
	}, nil
}

func (t *Tags) UpdateName(name string, newName string) (TagModel, error) {
	ctxSync, cancel := context.WithTimeout(context.Background(), dataCtxTimeAisle)
	defer cancel()

	res, err := t.neo.UpdateName(ctxSync, name, newName)
	if err != nil {
		return TagModel{}, err
	}

	return TagModel{
		Name:      res.Name,
		Status:    res.Status,
		Type:      res.Type,
		Color:     res.Color,
		Icon:      res.Icon,
		CreatedAt: res.CreatedAt,
	}, nil
}

func (t *Tags) Get(name string) (TagModel, error) {
	ctxSync, cancel := context.WithTimeout(context.Background(), dataCtxTimeAisle)
	defer cancel()

	neoTag, err := t.neo.Get(ctxSync, name)
	if err != nil {
		return TagModel{}, err
	}

	return TagModel{
		Name:      neoTag.Name,
		Status:    neoTag.Status,
		Type:      neoTag.Type,
		Color:     neoTag.Color,
		Icon:      neoTag.Icon,
		CreatedAt: neoTag.CreatedAt,
	}, nil
}

func (t *Tags) GetAll() ([]TagModel, error) {
	ctxSync, cancel := context.WithTimeout(context.Background(), dataCtxTimeAisle)
	defer cancel()

	neoTags, err := t.neo.GetAll(ctxSync)
	if err != nil {
		return nil, err
	}

	tags := make([]TagModel, len(neoTags))
	for i, neoTag := range neoTags {
		tags[i] = TagModel{
			Name:      neoTag.Name,
			Status:    neoTag.Status,
			Type:      neoTag.Type,
			Color:     neoTag.Color,
			Icon:      neoTag.Icon,
			CreatedAt: neoTag.CreatedAt,
		}
	}

	return tags, nil
}
