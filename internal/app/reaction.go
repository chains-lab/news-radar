package app

import (
	"context"

	"github.com/google/uuid"
)

func (a App) AddLike(ctx context.Context, userID, articleID uuid.UUID) error {
	return a.reactions.CreateLike(userID, articleID)
}

func (a App) RemoveLike(ctx context.Context, userID, articleID uuid.UUID) error {
	return a.reactions.RemoveLike(userID, articleID)
}

func (a App) AddRepost(ctx context.Context, userID, articleID uuid.UUID) error {
	return a.reactions.CreateRepost(userID, articleID)
}
