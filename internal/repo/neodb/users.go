package neodb

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/neo4j"
)

type UserModel struct {
	ID uuid.UUID `json:"id"`
}

type UsersImpl struct {
	driver neo4j.Driver
}

func NewUsers(uri, username, password string) (*UsersImpl, error) {
	driver, err := neo4j.NewDriver(uri, neo4j.BasicAuth(username, password, ""), func(c *neo4j.Config) {
		c.Encrypted = false
	})
	if err != nil {
		return nil, err
	}

	if err = driver.VerifyConnectivity(); err != nil {
		return nil, err
	}

	return &UsersImpl{
		driver: driver,
	}, nil
}

type UserCreateInput struct {
	ID uuid.UUID `json:"id"`
}

func (u *UsersImpl) Create(ctx context.Context, input UserCreateInput) error {
	session, err := u.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	if err != nil {
		return err
	}

	defer session.Close()

	resultChan := make(chan error, 1)

	go func() {
		_, err = session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
			cypher := `
				CREATE (u:User {id: $id})
			`
			_, err := tx.Run(cypher, map[string]interface{}{
				"id": input.ID.String(),
			})

			return nil, err
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

func (u *UsersImpl) Delete(ctx context.Context, id uuid.UUID) error {
	session, err := u.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	if err != nil {
		return err
	}

	defer session.Close()

	resultChan := make(chan error, 1)

	go func() {
		_, err = session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
			cypher := `
				MATCH (u:User {id: $id})
				DETACH DELETE u
			`
			_, err := tx.Run(cypher, map[string]interface{}{
				"id": id.String(),
			})

			return nil, err
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

func (u *UsersImpl) Get(ctx context.Context, id uuid.UUID) (UserModel, error) {
	session, err := u.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	if err != nil {
		return UserModel{}, err
	}
	defer session.Close()

	type resultWrapper struct {
		user UserModel
		err  error
	}
	resultChan := make(chan resultWrapper, 1)

	go func() {
		result, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
			cypher := `
				MATCH (u:User {id: $id})
				RETURN u
			`

			res, err := tx.Run(cypher, map[string]interface{}{
				"id": id.String(),
			})
			if err != nil {
				return nil, err
			}

			if !res.Next() {
				return nil, fmt.Errorf("user not found")
			}

			record := res.Record()
			// Приводим результат к neo4j.Node
			node, ok := record.GetByIndex(0).(neo4j.Node)
			if !ok {
				return nil, fmt.Errorf("unexpected type: %T", record.GetByIndex(0))
			}
			props := node.Props()
			// Извлекаем значение свойства "id"
			idProp, ok := props["id"].(string)
			if !ok {
				return nil, fmt.Errorf("id property missing or not a string")
			}

			user := UserModel{}
			user.ID, err = uuid.Parse(idProp)
			if err != nil {
				return nil, err
			}

			return user, nil
		})
		if err != nil {
			resultChan <- resultWrapper{UserModel{}, err}
			return
		}
		user, ok := result.(UserModel)
		if !ok {
			resultChan <- resultWrapper{UserModel{}, fmt.Errorf("unexpected result type")}
			return
		}
		resultChan <- resultWrapper{user, nil}
	}()

	select {
	case res := <-resultChan:
		return res.user, res.err
	case <-ctx.Done():
		return UserModel{}, ctx.Err()
	}
}
