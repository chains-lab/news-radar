package neodb

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hs-zavet/news-radar/internal/enums"
	"github.com/neo4j/neo4j-go-driver/neo4j"
)

//CREATE CONSTRAINT unique_tag_name IF NOT EXISTS
//FOR (t:Tag)
//REQUIRE t.name IS UNIQUE;

//CREATE CONSTRAINT unique_tag_id IF NOT EXISTS
//FOR (t:Tag)
//REQUIRE t.id IS UNIQUE;

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

func (t *TagsImpl) Create(ctx context.Context, input TagCreateInput) (TagModel, error) {
	session, err := t.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	if err != nil {
		return TagModel{}, err
	}

	defer session.Close()

	result, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		cypher := `
			CREATE (t:Tag {
				id: $id,
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
			"id":      input.Name,
			"name":    input.Name,
			"status":  input.Status,
			"type":    input.Type,
			"color":   input.Color,
			"icon":    input.Icon,
			"created": input.CreatedAt,
		}

		cursor, err := tx.Run(cypher, params)
		if err != nil {
			return nil, fmt.Errorf("failed to run update: %w", err)
		}
		if !cursor.Next() {
			return nil, fmt.Errorf("tag not found")
		}

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
	})

	if err != nil {
		return TagModel{}, err
	}

	tag, ok := result.(TagModel)
	if !ok {
		return TagModel{}, fmt.Errorf("unexpected result type")
	}

	return tag, nil
}

func (t *TagsImpl) Delete(ctx context.Context, id string) error {
	session, err := t.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	if err != nil {
		return err
	}
	defer session.Close()

	_, err = session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		cypher := `
			MATCH (t:Tag { id: $id })
			DETACH DELETE t
		`

		params := map[string]any{
			"id": id,
		}
		_, err := tx.Run(cypher, params)
		if err != nil {
			return nil, fmt.Errorf("failed to delete tag: %w", err)
		}

		return nil, nil
	})

	if err != nil {
		return fmt.Errorf("failed to delete tag: %w", err)
	}

	return nil
}

func (t *TagsImpl) Get(ctx context.Context, id string) (TagModel, error) {
	session, err := t.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	if err != nil {
		return TagModel{}, err
	}
	defer session.Close()

	result, err := session.ReadTransaction(func(tx neo4j.Transaction) (any, error) {
		cypher := `
			MATCH (t:Tag)
			WHERE t.id CONTAINS $id
			RETURN t
		`

		params := map[string]any{
			"id": id,
		}

		cursor, err := tx.Run(cypher, params)
		if err != nil {
			return nil, err
		}
		if cursor.Next() {
			nodeVal, ok := cursor.Record().Get("t")
			if !ok {
				return TagModel{}, fmt.Errorf("failed to find tag %q", id)
			}

			node, ok := nodeVal.(neo4j.Node)
			if !ok {
				return TagModel{}, fmt.Errorf("unexpected type for tag node")
			}

			props := node.Props()

			tag, err := parseTagFromProps(props)
			if err != nil {
				return TagModel{}, err
			}
			return tag, nil
		}
		return TagModel{}, fmt.Errorf("failed to find tag")
	})

	if err != nil {
		return TagModel{}, err
	}

	tag, ok := result.(TagModel)
	if !ok {
		return TagModel{}, fmt.Errorf("unexpected result type")
	}

	return tag, nil
}

func (t *TagsImpl) GetAll(ctx context.Context) ([]TagModel, error) {
	session, err := t.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	if err != nil {
		return nil, err
	}
	defer session.Close()

	result, err := session.ReadTransaction(func(tx neo4j.Transaction) (any, error) {
		cypher := `
            MATCH (t:Tag)
            OPTIONAL MATCH (t)<-[r:HAS_TAG]-(:Article)
            WITH t, count(r) AS popularity
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
		return nil, err
	}

	tags, ok := result.([]TagModel)
	if !ok {
		return nil, fmt.Errorf("unexpected result type")
	}

	return tags, nil
}

type TagUpdateInput struct {
	NewName   *string          `json:"name"`
	Status    *enums.TagStatus `json:"status"`
	Type      *enums.TagType   `json:"type"`
	Color     *string          `json:"color"`
	Icon      *string          `json:"icon"`
	UpdatedAt time.Time        `json:"updated_at,omitempty"`
}

func (t *TagsImpl) Update(ctx context.Context, id string, input TagUpdateInput) (TagModel, error) {
	setClauses := []string{"t.updated_at = $updated_at"}
	params := map[string]any{
		"current_id": id,
		"updated_at": input.UpdatedAt,
	}

	if input.NewName != nil {
		if *input.NewName == "" {
			return TagModel{}, fmt.Errorf("new name cannot be empty")
		}
		setClauses = append(setClauses, "t.id = $new_id")
		params["new_id"] = strings.ToLower(*input.NewName)
		setClauses = append(setClauses, "t.name = $new_name")
		params["new_name"] = *input.NewName
	}

	if input.Status != nil {
		if *input.Status == "" {
			return TagModel{}, fmt.Errorf("status cannot be empty")
		}
		setClauses = append(setClauses, "t.status = $status")
		params["status"] = string(*input.Status)
	}

	if input.Type != nil {
		if *input.Type == "" {
			return TagModel{}, fmt.Errorf("type cannot be empty")
		}
		setClauses = append(setClauses, "t.type = $type")
		params["type"] = string(*input.Type)
	}

	if input.Color != nil {
		if *input.Color == "" {
			return TagModel{}, fmt.Errorf("color cannot be empty")
		}
		setClauses = append(setClauses, "t.color = $color")
		params["color"] = *input.Color
	}

	if input.Icon != nil {
		if *input.Icon == "" {
			return TagModel{}, fmt.Errorf("icon cannot be empty")
		}
		setClauses = append(setClauses, "t.icon = $icon")
		params["icon"] = *input.Icon

	}

	if len(setClauses) == 1 {
		return t.Get(ctx, id)
	}

	session, err := t.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	if err != nil {
		return TagModel{}, err
	}
	defer session.Close()

	cypher := fmt.Sprintf(
		`MATCH (t:Tag { id: $current_id })
		 SET  %s
		 RETURN t`, strings.Join(setClauses, ", "))

	result, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		cursor, err := tx.Run(cypher, params)
		if err != nil {
			return TagModel{}, fmt.Errorf("failed to run update: %w", err)
		}
		if !cursor.Next() {
			return TagModel{}, fmt.Errorf("tag not found")
		}

		nodeVal, ok := cursor.Record().Get("t")
		if !ok {
			return TagModel{}, fmt.Errorf("failed to find tag")
		}

		node, ok := nodeVal.(neo4j.Node)
		if !ok {
			return TagModel{}, fmt.Errorf("unexpected type for tag node")
		}

		props := node.Props()

		tag, err := parseTagFromProps(props)
		if err != nil {
			return TagModel{}, err
		}

		return tag, nil
	})

	if err != nil {
		return TagModel{}, err
	}

	tag, ok := result.(TagModel)
	if !ok {
		return TagModel{}, fmt.Errorf("unexpected result type")
	}

	return tag, nil
}

func parseTagFromProps(props map[string]any) (TagModel, error) {
	var tag TagModel

	id, ok := props["id"].(string)
	if !ok || id == "" {
		return tag, fmt.Errorf("invalid or missing tag id")
	}

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

	color, ok := props["color"].(string)
	if !ok {
		return tag, fmt.Errorf("invalid or missing tag color")
	}

	icon, ok := props["icon"].(string)
	if !ok {
		return tag, fmt.Errorf("invalid or missing tag icon")
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
		ID:        id,
		Name:      name,
		Status:    status,
		Type:      tagType,
		Color:     color,
		Icon:      icon,
		CreatedAt: createdAt,
	}

	updatedAt, ok := props["updated_at"].(time.Time)
	if ok {
		tag.UpdatedAt = &updatedAt
	}

	return tag, nil
}
