package neo

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/neo4j"
	"github.com/recovery-flow/news-radar/internal/service/domain/models"
)

type Article struct {
	ID        uuid.UUID
	CreatedAt time.Time
	Tags      []string
	Themes    []string
}

type Articles interface {
	Create(ctx context.Context, article *Article) error
	Delete(ctx context.Context, ID uuid.UUID) error
	
	GetAllHasTag(ctx context.Context, ID uuid.UUID) ([]*models.Tag, error)
	SetHasTag(ctx context.Context, ID uuid.UUID, tags []string) error
	DeleteHasTag(ctx context.Context, ID uuid.UUID) error
	CreateHasTagRelationship(ctx context.Context, articleID uuid.UUID, tagName string) error
	DeleteHasTagRelationship(ctx context.Context, articleID uuid.UUID, tagName string) error

	GetAllAbout(ctx context.Context, ID uuid.UUID) ([]*models.Theme, error)
	SetAbout(ctx context.Context, ID uuid.UUID, themes []string) error
	DeleteAbout(ctx context.Context, ID uuid.UUID) error
	CreateAboutRelationship(ctx context.Context, articleID uuid.UUID, themeName string) error
	DeleteAboutRelationship(ctx context.Context, articleID uuid.UUID, themeName string) error

	GetAllAuthors(ctx context.Context, ID uuid.UUID) ([]uuid.UUID, error)
	SetAuthors(ctx context.Context, ID uuid.UUID, author []uuid.UUID) error
	DeleteAuthor(ctx context.Context, ID uuid.UUID) error
	CreateAuthorshipRelationship(ctx context.Context, articleID uuid.UUID, authorID uuid.UUID) error
	DeleteAuthorshipRelationship(ctx context.Context, articleID uuid.UUID, authorID uuid.UUID) error
}

type articles struct {
	driver neo4j.Driver
}

func NewArticles(uri, username, password string) (Articles, error) {
	driver, err := neo4j.NewDriver(uri, neo4j.BasicAuth(username, password, ""))
	if err != nil {
		return nil, fmt.Errorf("failed to create neo4j driver: %w", err)
	}

	if err = driver.VerifyConnectivity(); err != nil {
		return nil, fmt.Errorf("failed to verify connectivity: %w", err)
	}

	return &articles{
		driver: driver,
	}, nil
}

func (a *articles) Create(ctx context.Context, article *Article) error {
	if len(article.Tags) > 10 {
		return fmt.Errorf("article cannot have more than 10 tags")
	}
	if len(article.Themes) > 5 {
		return fmt.Errorf("article cannot have more than 5 themes")
	}

	session, err := a.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	if err != nil {
		return err
	}
	defer session.Close()

	_, err = session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		cypher := `
			CREATE (a:Article { id: $id, created_at: $created_at })
			FOREACH (tagName IN $tags |
				MATCH (t:Tag { name: tagName })
				MERGE (a)-[:HAS_TAG]->(t)
			)
			FOREACH (themeName IN $themes |
				MATCH (th:Theme { name: themeName })
				MERGE (a)-[:ABOUT]->(th)
			)
			RETURN a
		`
		params := map[string]any{
			"id":         article.ID.String(),
			"created_at": article.CreatedAt.UTC().Format(time.RFC3339),
			"tags":       article.Tags,
			"themes":     article.Themes,
		}

		_, err := tx.Run(cypher, params)
		if err != nil {
			return nil, fmt.Errorf("failed to create article with relationships: %w", err)
		}
		return nil, nil
	})

	return err
}

func (a *articles) Delete(ctx context.Context, id uuid.UUID) error {
	session, err := a.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	if err != nil {
		return err
	}
	defer session.Close()

	_, err = session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		cypher := `
			MATCH (a:Article { id: $id })
			DETACH DELETE a
		`
		params := map[string]any{
			"id": id.String(),
		}

		_, err := tx.Run(cypher, params)
		if err != nil {
			return nil, fmt.Errorf("failed to delete article: %w", err)
		}

		return nil, nil
	})

	return err
}

//Tags Relationship

func (a *articles) CreateHasTagRelationship(ctx context.Context, articleID uuid.UUID, tagName string) error {
	session, err := a.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	if err != nil {
		return err
	}
	defer session.Close()

	_, err = session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		cypher := `
			MATCH (art:Article { id: $articleID })
			MATCH (t:Tag { name: $tagName })
			MERGE (art)-[r:HAS_TAG]->(t)
		`
		params := map[string]any{
			"articleID": articleID.String(),
			"tagName":   tagName,
		}

		_, err := tx.Run(cypher, params)
		if err != nil {
			return nil, fmt.Errorf("failed to create HAS_TAG relationship: %w", err)
		}

		return nil, nil
	})

	return err
}

func (a *articles) DeleteHasTagRelationship(ctx context.Context, articleID uuid.UUID, tagName string) error {
	session, err := a.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	if err != nil {
		return err
	}
	defer session.Close()

	_, err = session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		cypher := `
			MATCH (art:Article { id: $articleID })-[r:HAS_TAG]->(t:Tag { name: $tagName })
			DELETE r
		`
		params := map[string]any{
			"articleID": articleID.String(),
			"tagName":   tagName,
		}

		_, err := tx.Run(cypher, params)
		if err != nil {
			return nil, fmt.Errorf("failed to delete HAS_TAG relationship: %w", err)
		}

		return nil, nil
	})

	return err
}

func (a *articles) GetAllHasTag(ctx context.Context, id uuid.UUID) ([]*models.Tag, error) {
	session, err := a.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	if err != nil {
		return nil, err
	}
	defer session.Close()

	result, err := session.ReadTransaction(func(tx neo4j.Transaction) (any, error) {
		cypher := `
			MATCH (a:Article { id: $id })-[:HAS_TAG]->(t:Tag)
			RETURN t
		`
		params := map[string]any{
			"id": id.String(),
		}

		records, err := tx.Run(cypher, params)
		if err != nil {
			return nil, err
		}

		var tagsList []*models.Tag
		for records.Next() {
			record := records.Record()
			node, ok := record.Get("t")
			if !ok {
				continue
			}
			props := node.(neo4j.Node).Props()
			tag := &models.Tag{
				Name:   props["name"].(string),
				Status: models.TagStatus(props["status"].(string)),
				Type:   models.TagType(props["category"].(string)),
			}
			if createdAtStr, ok := props["created_at"].(string); ok {
				if parsedTime, err := time.Parse(time.RFC3339, createdAtStr); err == nil {
					tag.CreatedAt = parsedTime
				}
			}
			tagsList = append(tagsList, tag)
		}
		return tagsList, nil
	})
	if err != nil {
		return nil, err
	}
	return result.([]*models.Tag), nil
}

func (a *articles) DeleteHasTag(ctx context.Context, id uuid.UUID) error {
	session, err := a.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	if err != nil {
		return err
	}
	defer session.Close()

	_, err = session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		cypher := `
			MATCH (a:Article { id: $id })-[r:HAS_TAG]->(:Tag)
			DELETE r
		`
		params := map[string]any{
			"id": id.String(),
		}
		_, err := tx.Run(cypher, params)
		if err != nil {
			return nil, fmt.Errorf("failed to delete HAS_TAG relationships: %w", err)
		}
		return nil, nil
	})
	return err
}

func (a *articles) SetHasTag(ctx context.Context, id uuid.UUID, tags []string) error {
	session, err := a.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	if err != nil {
		return err
	}
	defer session.Close()

	_, err = session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		deleteCypher := `
			MATCH (a:Article { id: $id })-[r:HAS_TAG]->(:Tag)
			DELETE r
		`
		params := map[string]any{"id": id.String()}
		_, err := tx.Run(deleteCypher, params)
		if err != nil {
			return nil, fmt.Errorf("failed to delete existing HAS_TAG relationships: %w", err)
		}

		createCypher := `
			MATCH (a:Article { id: $id })
			FOREACH (tagName IN $tags |
				MATCH (t:Tag { name: tagName })
				MERGE (a)-[:HAS_TAG]->(t)
			)
		`
		params["tags"] = tags
		_, err = tx.Run(createCypher, params)
		if err != nil {
			return nil, fmt.Errorf("failed to create new HAS_TAG relationships: %w", err)
		}
		return nil, nil
	})
	return err
}

//About relationships

func (a *articles) CreateAboutRelationship(ctx context.Context, articleID uuid.UUID, themeName string) error {
	session, err := a.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	if err != nil {
		return err
	}
	defer session.Close()

	_, err = session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		cypher := `
			MATCH (art:Article { id: $articleID })
			MATCH (th:Theme { name: $themeName })
			MERGE (art)-[r:ABOUT]->(th)
		`
		params := map[string]any{
			"articleID": articleID.String(),
			"themeName": themeName,
		}

		_, err := tx.Run(cypher, params)
		if err != nil {
			return nil, fmt.Errorf("failed to create ABOUT relationship: %w", err)
		}
		return nil, nil
	})

	return err
}

func (a *articles) DeleteAboutRelationship(ctx context.Context, articleID uuid.UUID, themeName string) error {
	session, err := a.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	if err != nil {
		return err
	}
	defer session.Close()

	_, err = session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		cypher := `
			MATCH (art:Article { id: $articleID })-[r:ABOUT]->(th:Theme { name: $themeName })
			DELETE r
		`
		params := map[string]any{
			"articleID": articleID.String(),
			"themeName": themeName,
		}

		_, err := tx.Run(cypher, params)
		if err != nil {
			return nil, fmt.Errorf("failed to delete ABOUT relationship: %w", err)
		}
		return nil, nil
	})

	return err
}

func (a *articles) GetAllAbout(ctx context.Context, id uuid.UUID) ([]*models.Theme, error) {
	session, err := a.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	if err != nil {
		return nil, err
	}
	defer session.Close()

	result, err := session.ReadTransaction(func(tx neo4j.Transaction) (any, error) {
		cypher := `
			MATCH (a:Article { id: $id })-[:ABOUT]->(th:Theme)
			RETURN th
		`
		params := map[string]any{
			"id": id.String(),
		}

		records, err := tx.Run(cypher, params)
		if err != nil {
			return nil, err
		}

		var themesList []*models.Theme
		for records.Next() {
			record := records.Record()
			node, ok := record.Get("th")
			if !ok {
				continue
			}
			props := node.(neo4j.Node).Props()
			theme := &models.Theme{
				Name:   props["name"].(string),
				Status: models.ThemeStatus(props["status"].(string)),
				Type:   models.ThemeType(props["category"].(string)),
			}
			if createdAtStr, ok := props["created_at"].(string); ok {
				if parsedTime, err := time.Parse(time.RFC3339, createdAtStr); err == nil {
					theme.CreatedAt = parsedTime
				}
			}
			themesList = append(themesList, theme)
		}
		return themesList, nil
	})
	if err != nil {
		return nil, err
	}
	return result.([]*models.Theme), nil
}

func (a *articles) DeleteAbout(ctx context.Context, id uuid.UUID) error {
	session, err := a.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	if err != nil {
		return err
	}
	defer session.Close()

	_, err = session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		cypher := `
			MATCH (a:Article { id: $id })-[r:ABOUT]->(:Theme)
			DELETE r
		`
		params := map[string]any{
			"id": id.String(),
		}
		_, err := tx.Run(cypher, params)
		if err != nil {
			return nil, fmt.Errorf("failed to delete ABOUT relationships: %w", err)
		}
		return nil, nil
	})
	return err
}

func (a *articles) SetAbout(ctx context.Context, id uuid.UUID, themes []string) error {
	session, err := a.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	if err != nil {
		return err
	}
	defer session.Close()

	_, err = session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		deleteCypher := `
			MATCH (a:Article { id: $id })-[r:ABOUT]->(:Theme)
			DELETE r
		`
		params := map[string]any{"id": id.String()}
		_, err := tx.Run(deleteCypher, params)
		if err != nil {
			return nil, fmt.Errorf("failed to delete existing ABOUT relationships: %w", err)
		}

		createCypher := `
			MATCH (a:Article { id: $id })
			FOREACH (themeName IN $themes |
				MATCH (th:Theme { name: themeName })
				MERGE (a)-[:ABOUT]->(th)
			)
		`
		params["themes"] = themes
		_, err = tx.Run(createCypher, params)
		if err != nil {
			return nil, fmt.Errorf("failed to create new ABOUT relationships: %w", err)
		}
		return nil, nil
	})
	return err
}

//Authorship relationships

func (a *articles) CreateAuthorshipRelationship(ctx context.Context, articleID uuid.UUID, authorID uuid.UUID) error {
	session, err := a.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	if err != nil {
		return err
	}
	defer session.Close()

	_, err = session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		cypher := `
            MATCH (art:Article { id: $articleID })
            MATCH (auth:Author { id: $authorID })
            MERGE (art)-[:AUTHORED_BY]->(auth)
        `
		params := map[string]any{
			"articleID": articleID.String(),
			"authorID":  authorID.String(),
		}
		_, err := tx.Run(cypher, params)
		if err != nil {
			return nil, fmt.Errorf("failed to create authorship relationship: %w", err)
		}
		return nil, nil
	})
	return err
}

func (a *articles) DeleteAuthorshipRelationship(ctx context.Context, articleID uuid.UUID, authorID uuid.UUID) error {
	session, err := a.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	if err != nil {
		return err
	}
	defer session.Close()

	_, err = session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		cypher := `
			MATCH (art:Article { id: $articleID })-[r:AUTHORED_BY]->(auth:Author { id: $authorID })
			DELETE r
		`
		params := map[string]any{
			"articleID": articleID.String(),
			"authorID":  authorID,
		}

		_, err := tx.Run(cypher, params)
		if err != nil {
			return nil, fmt.Errorf("failed to delete AUTHOR relationship: %w", err)
		}
		return nil, nil
	})

	return err
}

func (a *articles) GetAllAuthors(ctx context.Context, id uuid.UUID) ([]uuid.UUID, error) {
	session, err := a.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	if err != nil {
		return nil, err
	}
	defer session.Close()

	result, err := session.ReadTransaction(func(tx neo4j.Transaction) (any, error) {
		cypher := `
			MATCH (a:Article { id: $id })-[:AUTHORED_BY]->(au:Author)
			RETURN au.id AS authorID
		`
		params := map[string]any{
			"id": id.String(),
		}
		records, err := tx.Run(cypher, params)
		if err != nil {
			return nil, err
		}
		var authorIDs []uuid.UUID
		for records.Next() {
			record := records.Record()
			authorIDVal, ok := record.Get("authorID")
			if !ok {
				continue
			}

			if idStr, ok := authorIDVal.(string); ok {
				uid, err := uuid.Parse(idStr)
				if err != nil {
					continue
				}
				authorIDs = append(authorIDs, uid)
			}
		}
		return authorIDs, nil
	})
	if err != nil {
		return nil, err
	}
	return result.([]uuid.UUID), nil
}

func (a *articles) SetAuthors(ctx context.Context, id uuid.UUID, authors []uuid.UUID) error {
	session, err := a.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	if err != nil {
		return err
	}
	defer session.Close()

	_, err = session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		deleteCypher := `
			MATCH (a:Article { id: $id })-[r:AUTHORED_BY]->(:Author)
			DELETE r
		`
		params := map[string]any{
			"id": id.String(),
		}
		_, err := tx.Run(deleteCypher, params)
		if err != nil {
			return nil, fmt.Errorf("failed to delete existing authorship relationships: %w", err)
		}

		authorIDs := make([]string, len(authors))
		for i, authID := range authors {
			authorIDs[i] = authID.String()
		}

		createCypher := `
			MATCH (a:Article { id: $id })
			FOREACH (authorId IN $authors |
				MATCH (au:Author { id: authorId })
				MERGE (a)-[:AUTHORED_BY]->(au)
			)
		`
		params["authors"] = authorIDs

		_, err = tx.Run(createCypher, params)
		if err != nil {
			return nil, fmt.Errorf("failed to create new authorship relationships: %w", err)
		}

		return nil, nil
	})
	return err
}

func (a *articles) DeleteAuthor(ctx context.Context, id uuid.UUID) error {
	session, err := a.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	if err != nil {
		return err
	}
	defer session.Close()

	_, err = session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		cypher := `
			MATCH (a:Article { id: $id })-[r:AUTHORED_BY]->(:Author)
			DELETE r
		`
		params := map[string]any{
			"id": id.String(),
		}
		_, err := tx.Run(cypher, params)
		if err != nil {
			return nil, fmt.Errorf("failed to delete authorship relationships: %w", err)
		}
		return nil, nil
	})
	return err
}
