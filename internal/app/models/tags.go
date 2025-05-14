package models

import (
	"time"

	"github.com/chains-lab/news-radar/internal/enums"
)

type Tag struct {
	ID        string          `json:"id"`
	Name      string          `json:"name"`
	Status    enums.TagStatus `json:"status"`
	Type      enums.TagType   `json:"type"`
	Color     string          `json:"color"`
	Icon      string          `json:"icon"`
	UpdatedAt *time.Time      `json:"updated_at"`
	CreatedAt time.Time       `json:"created_at"`
}
