package repo

import (
	"context"

	"github.com/recovery-flow/news-radar/internal/app/models"
	"github.com/recovery-flow/news-radar/internal/config"
	"github.com/recovery-flow/news-radar/internal/data/neodb"
	"github.com/recovery-flow/news-radar/internal/data/redisdb"
)

type ThemesRedis interface {
	Add(ctx context.Context, tag redisdb.ThemeModels) error
	Get(ctx context.Context, tag string) (*redisdb.ThemeModels, error)
	Delete(ctx context.Context, tag string) error

	UpdateIcon(ctx context.Context, tag string, icon string) error
	UpdateColor(ctx context.Context, tag string, color string) error

	Drop(ctx context.Context) error
}

type ThemesNeo interface {
	Create(ctx context.Context, theme neodb.ThemeModels) error
	Delete(ctx context.Context, name string) error
	Get(ctx context.Context, name string) (*neodb.ThemeModels, error)

	UpdateStatus(ctx context.Context, name string, status models.ThemeStatus) error
	//UpdateName(ctx context.Context, name string, newName string) error

	Select(ctx context.Context) ([]neodb.ThemeModels, error)
}

type Themes struct {
	redis ThemesRedis
	neo   ThemesNeo
}

func NewThemes(cfg config.Config) (*Themes, error) {
	neo, err := neodb.NewThemes(cfg.Database.Neo4j.URI, cfg.Database.Neo4j.Username, cfg.Database.Neo4j.Password)
	if err != nil {
		return nil, err
	}
	redis := redisdb.NewThemes(cfg.Database.Redis.Addr, cfg.Database.Redis.Password, cfg.Database.Redis.DB)
	return &Themes{
		neo:   neo,
		redis: redis,
	}, nil
}

func (t *Themes) Create(ctx context.Context, theme models.Theme) error {
	neoTheme := neodb.ThemeModels{
		Name:      theme.Name,
		Status:    theme.Status,
		CreatedAt: theme.CreatedAt,
	}
	err := t.neo.Create(ctx, neoTheme)
	if err != nil {
		return err
	}

	err = t.redis.Add(ctx, redisdb.ThemeModels{
		Name:  theme.Name,
		Color: theme.Color,
		Icon:  theme.Icon,
	})
	if err != nil {
		return err
	}

	return nil
}

func (t *Themes) Delete(ctx context.Context, name string) error {
	err := t.neo.Delete(ctx, name)
	if err != nil {
		return err
	}

	err = t.redis.Delete(ctx, name)
	if err != nil {
		return err
	}

	return nil
}

func (t *Themes) Update(ctx context.Context, name string, fields map[string]any) error {
	if _, ok := fields["color"]; ok {
		err := t.redis.UpdateColor(ctx, name, fields["color"].(string))
		if err != nil {
			return err
		}
	}

	if _, ok := fields["icon"]; ok {
		err := t.redis.UpdateIcon(ctx, name, fields["icon"].(string))
		if err != nil {
			return err
		}
	}

	if _, ok := fields["status"]; ok {
		status, err := models.ParseThemeStatus(fields["status"].(string))
		if err != nil {
			return err
		}
		err = t.neo.UpdateStatus(ctx, name, status)
		if err != nil {
			return err
		}
	}

	return nil
}

func (t *Themes) Get(ctx context.Context, name string) (*models.Theme, error) {
	neoTheme, err := t.neo.Get(ctx, name)
	if err != nil {
		return nil, err
	}

	redisTheme, err := t.redis.Get(ctx, name)
	if err != nil {
		return nil, err
	}

	theme := createModelsTheme(*neoTheme, *redisTheme)

	return &theme, nil
}

func createModelsTheme(neo neodb.ThemeModels, redis redisdb.ThemeModels) models.Theme {
	return models.Theme{
		Name:      neo.Name,
		Color:     redis.Color,
		Icon:      redis.Icon,
		Status:    neo.Status,
		CreatedAt: neo.CreatedAt,
	}
}
