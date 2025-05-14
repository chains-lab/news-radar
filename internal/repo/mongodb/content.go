package mongodb

import (
	"context"
	"fmt"
	"time"

	"github.com/chains-lab/news-radar/internal/content"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (a *ArticlesQ) DeleteContentSection(
	ctx context.Context,
	index int,
	updatedAt time.Time,
) error {
	// 1. Загрузим текущий документ
	var article ArticleModel
	if err := a.collection.FindOne(ctx, a.filters).Decode(&article); err != nil {
		return err
	}

	// 2. Проверяем индекс на валидность
	if index < 0 || index >= len(article.Content) {
		return fmt.Errorf("invalid index %d for content length %d", index, len(article.Content))
	}

	// 3. Формируем новый слайс без удаляемой секции
	newContent := make([]content.Section, 0, len(article.Content)-1)
	newContent = append(newContent, article.Content[:index]...)
	newContent = append(newContent, article.Content[index+1:]...)

	// 4. Сдвигаем ID у всех секций, которые шли после удалённой
	for i := range newContent {
		if newContent[i].ID > index {
			newContent[i].ID-- // уменьшаем ID на 1
		}
	}

	// 5. Делаем апдейт в Mongo
	update := bson.M{
		"$set": bson.M{
			"content":    newContent,
			"updated_at": updatedAt,
		},
	}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var updated ArticleModel
	if err := a.collection.
		FindOneAndUpdate(ctx, a.filters, update, opts).
		Decode(&updated); err != nil {
		return fmt.Errorf("failed to delete content at index %d: %w", index, err)
	}

	return nil
}

func (a *ArticlesQ) UpdateContentSection(
	ctx context.Context,
	section content.Section,
	updatedAt time.Time,
) error {
	var article ArticleModel
	if err := a.collection.FindOne(ctx, a.filters).Decode(&article); err != nil {
		return err
	}

	if section.ID < 0 || section.ID > len(article.Content) {
		return fmt.Errorf("invalid index %d for content length %d", section.ID, len(article.Content))
	}

	if section.Media == nil && len(section.Text) == 0 && section.Audio == nil {
		return fmt.Errorf("section is empty")
	}

	var newContent []content.Section
	if section.ID < len(article.Content) {
		newContent = make([]content.Section, len(article.Content))
		copy(newContent, article.Content)
		newContent[section.ID] = section
	} else {
		newContent = make([]content.Section, 0, len(article.Content)+1)
		newContent = append(newContent, article.Content...)
		newContent = append(newContent, section)
	}

	update := bson.M{
		"$set": bson.M{
			"content":    newContent,
			"updated_at": updatedAt,
		},
	}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updated ArticleModel
	if err := a.collection.
		FindOneAndUpdate(ctx, a.filters, update, opts).
		Decode(&updated); err != nil {
		return fmt.Errorf("failed to update content at index %d: %w", section.ID, err)
	}

	return nil
}
