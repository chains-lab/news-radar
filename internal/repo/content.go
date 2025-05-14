package repo

import (
	"context"
	"time"

	"github.com/chains-lab/news-radar/internal/content"
	"github.com/google/uuid"
)

func (a *ArticlesRepo) UpdateContentSection(ID uuid.UUID, section content.Section) error {
	ctxSync, cancel := context.WithTimeout(context.Background(), dataCtxTimeAisle)
	defer cancel()

	updatedAt := time.Now().UTC()

	if section.Text == nil && section.Media == nil && section.Audio == nil {
		err := a.mongo.New().FilterID(ID).DeleteContentSection(ctxSync, section.ID, updatedAt)
		if err != nil {
			return err
		}
	} else {
		err := a.mongo.New().FilterID(ID).UpdateContentSection(ctxSync, section, updatedAt)
		if err != nil {
			return err
		}
	}

	return nil
}
