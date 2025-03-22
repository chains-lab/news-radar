package app

import "github.com/recovery-flow/news-radar/internal/service/app/entities"

type App interface {
}

type app struct {
	articles entities.Articles
	authors  entities.Authors
	tags     entities.Tags
	themes   entities.Themes
}
