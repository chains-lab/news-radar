package app

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hs-zavet/news-radar/internal/app/ape"
	"github.com/hs-zavet/news-radar/internal/app/models"
	"github.com/hs-zavet/news-radar/internal/content"
	"github.com/hs-zavet/news-radar/internal/enums"
	"github.com/hs-zavet/news-radar/internal/repo"
	"go.mongodb.org/mongo-driver/mongo"
)

type CreateArticleRequest struct {
	Title string `json:"title"`
}

func (a App) CreateArticle(ctx context.Context, request CreateArticleRequest) (models.Article, error) {
	ArticleID := uuid.New()
	CreatedAt := time.Now().UTC()

	article, err := a.articles.Create(repo.ArticleCreateInput{
		ID:        ArticleID,
		Title:     request.Title,
		Status:    enums.ArticleStatusInactive,
		CreatedAt: CreatedAt,
	})
	if err != nil {
		return models.Article{}, err
	}

	return models.Article{
		ID:        article.ID,
		Title:     article.Title,
		Icon:      nil,
		Desc:      nil,
		Content:   nil,
		Status:    article.Status,
		UpdatedAt: nil,
		CreatedAt: article.CreatedAt,
	}, nil
}

type UpdateArticleRequest struct {
	Title  *string              `json:"title,omitempty"`
	Status *enums.ArticleStatus `json:"status,omitempty"`
	Icon   *string              `json:"icon,omitempty"`
	Desc   *string              `json:"desc,omitempty"`
}

func (a App) UpdateArticle(ctx context.Context, articleID uuid.UUID, request UpdateArticleRequest) (models.Article, error) {
	var err error

	input := repo.ArticleUpdateInput{}
	updated := false

	curArticle, err := a.articles.GetByID(articleID)
	if err != nil {
		switch {
		case errors.Is(err, mongo.ErrNoDocuments):
			return models.Article{}, ape.ErrArticleNotFound
		default:
			return models.Article{}, err
		}
	}

	if request.Title != nil {
		input.Title = request.Title
		updated = true
	}

	if request.Status != nil {
		input.Status = request.Status
		if *request.Status == enums.ArticleStatusPublished {
			if curArticle.PublishedAt == nil {
				publishedAt := time.Now().UTC()
				input.PublishedAt = &publishedAt
			}
		}
		updated = true
	}

	if request.Icon != nil {
		input.Icon = request.Icon
		updated = true
	}

	if request.Desc != nil {
		input.Desc = request.Desc
		updated = true
	}

	if !updated {
		//for idempotency
		return a.GetArticleByID(ctx, articleID)
	}

	article, err := a.articles.Update(articleID, input)
	if err != nil {
		return models.Article{}, err
	}

	res := ArticleRepoToModels(article)
	return res, nil
}

func (a App) UpdateArticleContent(ctx context.Context, articleID uuid.UUID, index int, section content.Section) (models.Article, error) {
	article, err := a.articles.UpdateContent(articleID, index, section)
	if err != nil {
		return models.Article{}, err
	}

	res := ArticleRepoToModels(article)

	return res, nil
}

func (a App) DeleteArticle(ctx context.Context, articleID uuid.UUID) error {
	_, err := a.articles.GetByID(articleID)
	if err != nil {
		return ape.ErrArticleNotFound
	}

	err = a.articles.Delete(articleID)
	if err != nil {
		return err
	}

	return nil
}

func (a App) GetArticleByID(ctx context.Context, articleID uuid.UUID) (models.Article, error) {
	article, err := a.articles.GetByID(articleID)
	if err != nil {
		return models.Article{}, ape.ErrArticleNotFound
	}

	res := ArticleRepoToModels(article)

	return res, nil
}

//HASHTAGS

func (a App) SetArticleTags(ctx context.Context, articleID uuid.UUID, tags []string) error {
	_, err := a.articles.GetByID(articleID)
	if err != nil {
		switch {
		case errors.Is(err, mongo.ErrNoDocuments):
			return ape.ErrArticleNotFound
		default:
			return err
		}
	}

	seen := make(map[string]struct{}, len(tags))
	for _, id := range tags {
		if _, exists := seen[id]; exists {
			return ape.ErrTagReplication
		}
		seen[id] = struct{}{}
	}

	if len(tags) > 10 {
		return ape.ErrTooManyTags
	}

	for _, tag := range tags {
		tagModel, err := a.tags.Get(tag)
		if err != nil {
			switch {
			case errors.Is(err, mongo.ErrNoDocuments):
				return ape.ErrTagNotFound
			default:
				return err
			}
		}

		if tagModel.Status != enums.TagStatusActive {
			return ape.ErrTagInactive
		}
	}

	err = a.articles.SetTags(articleID, tags)
	if err != nil {
		return err
	}

	return nil
}

func (a App) GetArticleForTags(ctx context.Context, tag string) ([]models.Article, error) {
	articles, err := a.articles.GetArticlesForTag(tag)
	if err != nil {
		return nil, err
	}

	var res []models.Article
	for _, articleID := range articles {
		article, err := a.articles.GetByID(articleID)
		if err != nil {
			return nil, err
		}

		elem := models.Article{
			ID:        article.ID,
			Status:    article.Status,
			Title:     article.Title,
			CreatedAt: article.CreatedAt,
		}

		if article.Desc != nil {
			elem.Desc = article.Desc
		}

		if article.Icon != nil {
			elem.Icon = article.Icon
		}

		if article.Content != nil {
			elem.Content = article.Content
		}

		if article.UpdatedAt != nil {
			elem.UpdatedAt = article.UpdatedAt
		}

		res = append(res, elem)
	}

	return res, nil
}

func (a App) GetArticleTags(ctx context.Context, articleID uuid.UUID) ([]models.Tag, error) {
	_, err := a.articles.GetByID(articleID)
	if err != nil {
		switch {
		case errors.Is(err, mongo.ErrNoDocuments):
			return nil, ape.ErrArticleNotFound
		default:
			return nil, err
		}
	}

	tags, err := a.articles.GetTags(articleID)
	if err != nil {
		return nil, err
	}

	var res []models.Tag
	for _, tagID := range tags {
		tag, err := a.tags.Get(tagID)
		if err != nil {
			return nil, err
		}

		res = append(res, models.Tag{
			Name:      tag.Name,
			Status:    tag.Status,
			Type:      tag.Type,
			Color:     tag.Color,
			Icon:      tag.Icon,
			CreatedAt: tag.CreatedAt,
		})
	}

	return res, nil
}

func (a App) AddArticleTag(ctx context.Context, articleID uuid.UUID, tagId string) error {
	_, err := a.articles.GetByID(articleID)
	if err != nil {
		switch {
		case errors.Is(err, mongo.ErrNoDocuments):
			return ape.ErrArticleNotFound
		default:
			return err
		}
	}

	tag, err := a.tags.Get(tagId)
	if err != nil {
		return ape.ErrTagNotFound
	}

	if tag.Status != enums.TagStatusActive {
		return ape.ErrTagInactive
	}

	err = a.articles.AddTag(articleID, tagId)
	if err != nil {
		return err
	}

	return nil
}

func (a App) DeleteArticleTag(ctx context.Context, articleID uuid.UUID, tagId string) error {
	_, err := a.articles.GetByID(articleID)
	if err != nil {
		switch {
		case errors.Is(err, mongo.ErrNoDocuments):
			return ape.ErrArticleNotFound
		default:
			return err
		}
	}

	_, err = a.tags.Get(tagId)
	if err != nil {
		return ape.ErrTagNotFound
	}

	err = a.articles.DeleteTag(articleID, strings.ToLower(tagId))
	if err != nil {
		return err
	}

	return nil
}

func (a App) CleanArticleTags(ctx context.Context, articleID uuid.UUID) error {
	_, err := a.articles.GetByID(articleID)
	if err != nil {
		switch {
		case errors.Is(err, mongo.ErrNoDocuments):
			return ape.ErrArticleNotFound
		default:
			return err
		}
	}

	tags, err := a.articles.GetTags(articleID)
	if err != nil {
		return err
	}

	for _, tag := range tags {
		err := a.articles.DeleteTag(articleID, tag)
		if err != nil {
			return err
		}
	}

	return nil
}

//AUTHORSHIP

func (a App) SetAuthors(ctx context.Context, articleID uuid.UUID, authors []uuid.UUID) error {
	_, err := a.articles.GetByID(articleID)
	if err != nil {
		switch {
		case errors.Is(err, mongo.ErrNoDocuments):
			return ape.ErrArticleNotFound
		default:
			return err
		}
	}

	seen := make(map[uuid.UUID]struct{}, len(authors))
	for _, id := range authors {
		if _, exists := seen[id]; exists {
			return ape.ErrAuthorReplication
		}
		seen[id] = struct{}{}
	}

	for _, author := range authors {
		authorModel, err := a.authors.GetByID(author)
		if err != nil {
			switch {
			case errors.Is(err, mongo.ErrNoDocuments):
				return ape.ErrAuthorNotFound
			default:
				return err
			}
		}

		if authorModel.Status != enums.AuthorStatusActive {
			return ape.ErrAuthorInactive
		}
	}

	err = a.articles.SetAuthors(articleID, authors)
	if err != nil {
		return err
	}

	return nil
}

func (a App) GetArticleAuthors(ctx context.Context, articleID uuid.UUID) ([]models.Author, error) {
	_, err := a.articles.GetByID(articleID)
	if err != nil {
		switch {
		case errors.Is(err, mongo.ErrNoDocuments):
			return nil, ape.ErrArticleNotFound
		default:
			return nil, err
		}
	}

	authors, err := a.articles.GetAuthors(articleID)
	if err != nil {
		return nil, err
	}

	var res []models.Author
	for _, authorID := range authors {
		author, err := a.authors.GetByID(authorID)
		if err != nil {
			return nil, err
		}

		res = append(res, models.Author{
			ID:        author.ID,
			Name:      author.Name,
			Status:    author.Status,
			Desc:      author.Desc,
			Avatar:    author.Avatar,
			Email:     author.Email,
			Telegram:  author.Telegram,
			Twitter:   author.Twitter,
			CreatedAt: author.CreatedAt,
			UpdatedAt: author.UpdatedAt,
		})
	}

	return res, nil
}

func (a App) GetArticleForAuthor(ctx context.Context, authorID uuid.UUID) ([]models.Article, error) {
	_, err := a.authors.GetByID(authorID)
	if err != nil {
		switch {
		case errors.Is(err, mongo.ErrNoDocuments):
			return nil, ape.ErrAuthorNotFound
		default:
			return nil, err
		}
	}

	articles, err := a.articles.GetArticlesForAuthor(authorID)
	if err != nil {
		return nil, err
	}

	var res []models.Article
	for _, articleID := range articles {
		article, err := a.articles.GetByID(articleID)
		if err != nil {
			return nil, err
		}

		elem := models.Article{
			ID:        article.ID,
			Status:    article.Status,
			Title:     article.Title,
			CreatedAt: article.CreatedAt,
		}

		if article.Desc != nil {
			elem.Desc = article.Desc
		}

		if article.Icon != nil {
			elem.Icon = article.Icon
		}

		if article.Content != nil {
			elem.Content = article.Content
		}

		if article.UpdatedAt != nil {
			elem.UpdatedAt = article.UpdatedAt
		}

		if article.PublishedAt != nil {
			elem.PublishedAt = article.PublishedAt
		}

		res = append(res, elem)
	}

	return res, nil
}

func (a App) AddArticleAuthor(ctx context.Context, articleID uuid.UUID, authorID uuid.UUID) error {
	_, err := a.articles.GetByID(articleID)
	if err != nil {
		switch {
		case errors.Is(err, mongo.ErrNoDocuments):
			return ape.ErrArticleNotFound
		default:
			return err
		}
	}

	author, err := a.authors.GetByID(authorID)
	if err != nil {
		switch {
		case errors.Is(err, mongo.ErrNoDocuments):
			return ape.ErrAuthorNotFound
		default:
			return err
		}
	}

	if author.Status != enums.AuthorStatusActive {
		return ape.ErrAuthorInactive
	}

	err = a.articles.AddAuthor(articleID, authorID)
	if err != nil {
		return err
	}

	return nil
}

func (a App) DeleteArticleAuthor(ctx context.Context, articleID uuid.UUID, authorID uuid.UUID) error {
	_, err := a.articles.GetByID(articleID)
	if err != nil {
		switch {
		case errors.Is(err, mongo.ErrNoDocuments):
			return ape.ErrArticleNotFound
		default:
			return err
		}
	}

	err = a.articles.DeleteAuthor(articleID, authorID)
	if err != nil {
		return err
	}

	return nil
}

func (a App) CleanArticleAuthors(ctx context.Context, articleID uuid.UUID) error {
	_, err := a.articles.GetByID(articleID)
	if err != nil {
		switch {
		case errors.Is(err, mongo.ErrNoDocuments):
			return ape.ErrArticleNotFound
		default:
			return err
		}
	}

	authors, err := a.articles.GetAuthors(articleID)
	if err != nil {
		return err
	}

	for _, author := range authors {
		err := a.articles.DeleteAuthor(articleID, author)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a App) RecommendByTopic(ctx context.Context, articleID uuid.UUID, limit int) ([]models.Article, error) {
	articles, err := a.articles.RecommendByTopic(ctx, articleID, limit)
	if err != nil {
		return nil, err
	}

	var res []models.Article
	for _, article := range articles {
		res = append(res, ArticleRepoToModels(article))
	}

	return res, nil
}

func (a App) TopicSearch(ctx context.Context, tag string, start, limit int) ([]models.Article, error) {
	articles, err := a.articles.TopicSearch(ctx, tag, start, limit)
	if err != nil {
		return nil, err
	}

	var res []models.Article
	for _, article := range articles {
		res = append(res, ArticleRepoToModels(article))
	}

	return res, nil
}

func ArticleRepoToModels(article repo.ArticleModel) models.Article {
	res := models.Article{
		ID:        article.ID,
		Title:     article.Title,
		Status:    article.Status,
		CreatedAt: article.CreatedAt,
	}
	if article.Desc != nil {
		res.Desc = article.Desc
	}
	if article.Icon != nil {
		res.Icon = article.Icon
	}
	if article.Content != nil {
		res.Content = article.Content
	}
	if article.UpdatedAt != nil {
		res.UpdatedAt = article.UpdatedAt
	}
	if article.PublishedAt != nil {
		res.PublishedAt = article.PublishedAt
	}

	return res
}
