package neodb

import (
	"context"
	"fmt"
	"time"

	"github.com/hs-zavet/news-radar/internal/enums"
	"github.com/neo4j/neo4j-go-driver/neo4j"
)

type TagModel struct {
	Name      string          `json:"name"`
	Status    enums.TagStatus `json:"status"`
	Type      enums.TagType   `json:"type"`
	Color     string          `json:"color"`
	Icon      string          `json:"icon"`
	CreatedAt time.Time       `json:"created_at"`
}

type TagsImpl struct {
	driver neo4j.Driver
}

func NewTags(uri, username, password string) (*TagsImpl, error) {
	driver, err := neo4j.NewDriver(uri, neo4j.BasicAuth(username, password, ""), func(c *neo4j.Config) {
		c.Encrypted = false
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create neo4j driver: %w", err)
	}

	if err = driver.VerifyConnectivity(); err != nil {
		return nil, fmt.Errorf("failed to verify connectivity: %w", err)
	}

	return &TagsImpl{
		driver: driver,
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

func (t *TagsImpl) Create(ctx context.Context, input TagCreateInput) error {
	session, err := t.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	if err != nil {
		return err
	}

	defer session.Close()

	resultChan := make(chan error, 1)

	go func() {
		_, err = session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
			cypher := `
				CREATE (t:Tag {
					name: $name,
					status: $status,
					type: $type,
					color: $color,
					icon: $icon,
					created_at: $created
				})
				RETURN t
			`

			params := map[string]any{
				"name":    input.Name,
				"status":  input.Status,
				"type":    input.Type,
				"color":   input.Color,
				"icon":    input.Icon,
				"created": input.CreatedAt,
			}

			_, err := tx.Run(cypher, params)
			if err != nil {
				return nil, fmt.Errorf("failed to create tag: %w", err)
			}

			return nil, nil
		})
		resultChan <- err
	}()

	select {
	case err := <-resultChan:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (t *TagsImpl) Delete(ctx context.Context, name string) error {
	session, err := t.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	if err != nil {
		return err
	}

	defer session.Close()

	resultChan := make(chan error, 1)

	go func() {
		_, err = session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
			cypher := `
				MATCH (t:Tag { name: $name })
				DETACH DELETE t
			`

			params := map[string]any{
				"name": name,
			}
			_, err := tx.Run(cypher, params)
			if err != nil {
				return nil, fmt.Errorf("failed to delete tag: %w", err)
			}

			return nil, nil
		})
		resultChan <- err
	}()

	select {
	case err := <-resultChan:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (t *TagsImpl) UpdateStatus(ctx context.Context, name string, status enums.TagStatus) error {
	session, err := t.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	if err != nil {
		return err
	}
	defer session.Close()

	errChan := make(chan error, 1)
	go func() {
		_, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
			cypher := `
				MATCH (t:Tag { name: $name })
				SET t.status = $status
				RETURN t
			`

			params := map[string]any{
				"name":   name,
				"status": string(status),
			}
			_, err := tx.Run(cypher, params)
			if err != nil {
				return nil, fmt.Errorf("failed to update tag status: %w", err)
			}
			return nil, nil
		})
		errChan <- err
	}()

	select {
	case err := <-errChan:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (t *TagsImpl) UpdateName(ctx context.Context, name string, newName string) error {
	session, err := t.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	if err != nil {
		return err
	}
	defer session.Close()

	errChan := make(chan error, 1)
	go func() {
		_, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
			cypher := `
				MATCH (t:Tag { name: $name })
				SET t.name = $newName
				RETURN t
			`

			params := map[string]any{
				"name":    name,
				"newName": newName,
			}
			_, err := tx.Run(cypher, params)
			if err != nil {
				return nil, fmt.Errorf("failed to update tag name: %w", err)
			}
			return nil, nil
		})
		errChan <- err
	}()

	select {
	case err := <-errChan:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (t *TagsImpl) UpdateType(ctx context.Context, name string, newType enums.TagType) error {
	session, err := t.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	if err != nil {
		return err
	}
	defer session.Close()

	errChan := make(chan error, 1)
	go func() {
		_, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
			cypher := `
				MATCH (t:Tag { name: $name })
				SET t.type = $newType
				RETURN t
			`

			params := map[string]any{
				"name":    name,
				"newType": string(newType),
			}
			_, err := tx.Run(cypher, params)
			if err != nil {
				return nil, fmt.Errorf("failed to update tag type: %w", err)
			}
			return nil, nil
		})
		errChan <- err
	}()

	select {
	case err := <-errChan:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (t *TagsImpl) UpdateColor(ctx context.Context, name string, newColor string) error {
	session, err := t.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	if err != nil {
		return err
	}
	defer session.Close()

	errChan := make(chan error, 1)
	go func() {
		_, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
			cypher := `
				MATCH (t:Tag { name: $name })
				SET t.color = $newColor
				RETURN t
			`

			params := map[string]any{
				"name":     name,
				"newColor": newColor,
			}
			_, err := tx.Run(cypher, params)
			if err != nil {
				return nil, fmt.Errorf("failed to update tag color: %w", err)
			}
			return nil, nil
		})
		errChan <- err
	}()

	select {
	case err := <-errChan:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (t *TagsImpl) UpdateIcon(ctx context.Context, name string, newIcon string) error {
	session, err := t.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	if err != nil {
		return err
	}
	defer session.Close()

	errChan := make(chan error, 1)
	go func() {
		_, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
			cypher := `
				MATCH (t:Tag { name: $name })
				SET t.icon = $newIcon
				RETURN t
			`

			params := map[string]any{
				"name":    name,
				"newIcon": newIcon,
			}
			_, err := tx.Run(cypher, params)
			if err != nil {
				return nil, fmt.Errorf("failed to update tag icon: %w", err)
			}
			return nil, nil
		})
		errChan <- err
	}()

	select {
	case err := <-errChan:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (t *TagsImpl) Get(ctx context.Context, name string) (TagModel, error) {
	session, err := t.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	if err != nil {
		return TagModel{}, err
	}
	defer session.Close()

	type resultWrapper struct {
		tag TagModel
		err error
	}
	resultChan := make(chan resultWrapper, 1)

	go func() {
		result, err := session.ReadTransaction(func(tx neo4j.Transaction) (any, error) {
			cypher := `
				MATCH (t:Tag)
				WHERE toLower(t.name) CONTAINS toLower($name)
				RETURN t
			`

			params := map[string]any{
				"name": name,
			}

			cursor, err := tx.Run(cypher, params)
			if err != nil {
				return nil, err
			}
			if cursor.Next() {
				nodeVal, ok := cursor.Record().Get("t")
				if !ok {
					return nil, fmt.Errorf("failed to find tag")
				}
				node, ok := nodeVal.(neo4j.Node)
				if !ok {
					return nil, fmt.Errorf("unexpected type for tag node")
				}

				props := node.Props()
				tag, err := parseTagFromProps(props)
				if err != nil {
					return nil, err
				}
				return tag, nil
			}
			return TagModel{}, fmt.Errorf("failed to find tag")
		})
		if err != nil {
			resultChan <- resultWrapper{TagModel{}, err}
			return
		}
		tag, ok := result.(TagModel)
		if !ok {
			resultChan <- resultWrapper{TagModel{}, fmt.Errorf("unexpected result type")}
			return
		}
		resultChan <- resultWrapper{tag, nil}
	}()

	select {
	case res := <-resultChan:
		return res.tag, res.err
	case <-ctx.Done():
		return TagModel{}, ctx.Err()
	}
}

func (t *TagsImpl) GetAll(ctx context.Context) ([]TagModel, error) {
	session, err := t.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	if err != nil {
		return nil, err
	}
	defer session.Close()

	type resultWrapper struct {
		tags []TagModel
		err  error
	}
	resultChan := make(chan resultWrapper, 1)

	go func() {
		result, err := session.ReadTransaction(func(tx neo4j.Transaction) (any, error) {
			cypher := `
				MATCH (t:Tag)
				OPTIONAL MATCH (t)<-[r:ABOUT]-(:Article)
				WITH t, count(r) as popularity
				RETURN t ORDER BY popularity DESC
			`

			cursor, err := tx.Run(cypher, nil)
			if err != nil {
				return nil, err
			}

			var tagsList []TagModel
			for cursor.Next() {
				record := cursor.Record()
				nodeVal, ok := record.Get("t")
				if !ok {
					continue
				}
				node, ok := nodeVal.(neo4j.Node)
				if !ok {
					continue
				}
				props := node.Props()
				tag, err := parseTagFromProps(props)
				if err != nil {

					continue
				}
				tagsList = append(tagsList, tag)
			}
			return tagsList, nil
		})
		if err != nil {
			resultChan <- resultWrapper{nil, err}
			return
		}
		tags, ok := result.([]TagModel)
		if !ok {
			resultChan <- resultWrapper{nil, fmt.Errorf("unexpected result type")}
			return
		}
		resultChan <- resultWrapper{tags, nil}
	}()

	select {
	case res := <-resultChan:
		return res.tags, res.err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func parseTagFromProps(props map[string]any) (TagModel, error) {
	var tag TagModel

	name, ok := props["name"].(string)
	if !ok || name == "" {
		return tag, fmt.Errorf("invalid or missing tag name")
	}

	statusStr, ok := props["status"].(string)
	if !ok || statusStr == "" {
		return tag, fmt.Errorf("invalid or missing tag status")
	}

	tagTypeStr, ok := props["type"].(string)
	if !ok || tagTypeStr == "" {
		return tag, fmt.Errorf("invalid or missing tag type")
	}

	status, ok := enums.ParseTagStatus(statusStr)
	if !ok {
		return tag, fmt.Errorf("invalid tag status value")
	}

	tagType, ok := enums.ParseTagType(tagTypeStr)
	if !ok {
		return tag, fmt.Errorf("invalid tag type value")
	}

	createdAt, ok := props["created_at"].(time.Time)
	if !ok {
		return tag, fmt.Errorf("invalid or missing created_at timestamp")
	}

	tag = TagModel{
		Name:      name,
		Status:    status,
		Type:      tagType,
		CreatedAt: createdAt,
	}
	return tag, nil
}
