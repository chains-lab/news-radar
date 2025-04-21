package repo

import (
	"context"
	"time"

	"github.com/hs-zavet/news-radar/internal/config"
	"github.com/hs-zavet/news-radar/internal/enums"
	"github.com/hs-zavet/news-radar/internal/repo/neodb"
)

type TagModel struct {
	ID        string          `json:"id"`
	Name      string          `json:"name"`
	Status    enums.TagStatus `json:"status"`
	Type      enums.TagType   `json:"type"`
	Color     string          `json:"color"`
	Icon      string          `json:"icon"`
	UpdatedAt *time.Time      `json:"updated_at,omitempty"`
	CreatedAt time.Time       `json:"created_at"`
}

type tagsNeo interface {
	Create(ctx context.Context, input neodb.TagCreateInput) (neodb.TagModel, error)
	Delete(ctx context.Context, id string) error

	Update(ctx context.Context, id string, input neodb.TagUpdateInput) (neodb.TagModel, error)

	Get(ctx context.Context, id string) (neodb.TagModel, error)
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

func (t *Tags) Create(input TagCreateInput) (TagModel, error) {
	ctxSync, cancel := context.WithTimeout(context.Background(), dataCtxTimeAisle)
	defer cancel()

	tag, err := t.neo.Create(ctxSync, neodb.TagCreateInput{
		Name:      input.Name,
		Status:    input.Status,
		Type:      input.Type,
		Color:     input.Color,
		Icon:      input.Icon,
		CreatedAt: input.CreatedAt,
	})
	if err != nil {
		return TagModel{}, err
	}

	return TagModel{
		ID:        tag.ID,
		Name:      tag.Name,
		Status:    tag.Status,
		Type:      tag.Type,
		Color:     tag.Color,
		Icon:      tag.Icon,
		CreatedAt: tag.CreatedAt,
	}, nil
}

func (t *Tags) Delete(id string) error {
	ctxSync, cancel := context.WithTimeout(context.Background(), dataCtxTimeAisle)
	defer cancel()

	err := t.neo.Delete(ctxSync, id)
	if err != nil {
		return err
	}

	return nil
}

type TagUpdateInput struct {
	Name   *string          `json:"name"`
	Status *enums.TagStatus `json:"status"`
	Type   *enums.TagType   `json:"type"`
	Color  *string          `json:"color"`
	Icon   *string          `json:"icon"`
}

func (t *Tags) Update(id string, input TagUpdateInput) (TagModel, error) {
	ctxSync, cancel := context.WithTimeout(context.Background(), dataCtxTimeAisle)
	defer cancel()

	updatedAt := time.Now().UTC()

	var neoInput neodb.TagUpdateInput

	if input.Name != nil {
		neoInput.NewName = input.Name
	}
	if input.Status != nil {
		neoInput.Status = input.Status
	}
	if input.Type != nil {
		neoInput.Type = input.Type
	}
	if input.Color != nil {
		neoInput.Color = input.Color
	}
	if input.Icon != nil {
		neoInput.Icon = input.Icon
	}
	if input.Status != nil {
		neoInput.Status = input.Status
	}
	neoInput.UpdatedAt = updatedAt

	neoTag, err := t.neo.Update(ctxSync, id, neoInput)
	if err != nil {
		return TagModel{}, err
	}

	return TagModel{
		ID:        neoTag.ID,
		Name:      neoTag.Name,
		Status:    neoTag.Status,
		Type:      neoTag.Type,
		Color:     neoTag.Color,
		Icon:      neoTag.Icon,
		CreatedAt: neoTag.CreatedAt,
		UpdatedAt: neoTag.UpdatedAt,
	}, nil
}

func (t *Tags) Get(id string) (TagModel, error) {
	ctxSync, cancel := context.WithTimeout(context.Background(), dataCtxTimeAisle)
	defer cancel()

	neoTag, err := t.neo.Get(ctxSync, id)
	if err != nil {
		return TagModel{}, err
	}

	tag := TagModel{
		ID:        neoTag.ID,
		Name:      neoTag.Name,
		Status:    neoTag.Status,
		Type:      neoTag.Type,
		Color:     neoTag.Color,
		Icon:      neoTag.Icon,
		CreatedAt: neoTag.CreatedAt,
	}
	if neoTag.UpdatedAt != nil {
		tag.UpdatedAt = neoTag.UpdatedAt
	}

	return tag, nil
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
			ID:        neoTag.ID,
			Name:      neoTag.Name,
			Status:    neoTag.Status,
			Type:      neoTag.Type,
			Color:     neoTag.Color,
			Icon:      neoTag.Icon,
			CreatedAt: neoTag.CreatedAt,
		}
		if neoTag.UpdatedAt != nil {
			tags[i].UpdatedAt = neoTag.UpdatedAt
		}
	}

	return tags, nil
}
