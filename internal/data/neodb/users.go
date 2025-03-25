package neodb

import (
	"context"

	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/neo4j"
)

type User struct {
	ID uuid.UUID `json:"id"`
}

type Users interface {
	Create(ctx context.Context, user User) error
	Get(ctx context.Context, id uuid.UUID) (*User, error)
}

type users struct {
	driver neo4j.Driver
}

func NewUsers(uri, username, password string) (Users, error) {
	driver, err := neo4j.NewDriver(uri, neo4j.BasicAuth(username, password, ""))
	if err != nil {
		return nil, err
	}

	if err = driver.VerifyConnectivity(); err != nil {
		return nil, err
	}

	return &users{
		driver: driver,
	}, nil
}

func (u *users) Create(ctx context.Context, user User) error {
	session, err := u.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	if err != nil {
		return err
	}
	defer session.Close()

	_, err = session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		cypher := `
			CREATE (u:User {id: $id})
		`
		_, err := tx.Run(cypher, map[string]interface{}{
			"id": user.ID.String(),
		})
		return nil, err
	})
	return err
}

func (u *users) Delete(ctx context.Context, id uuid.UUID) error {
	session, err := u.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	if err != nil {
		return err
	}
	defer session.Close()

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
	return err
}

func (u *users) Get(ctx context.Context, id uuid.UUID) (*User, error) {
	session, err := u.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	if err != nil {
		return nil, err
	}
	defer session.Close()

	result, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		cypher := `
			MATCH (u:User {id: $id})
			RETURN u
		`
		result, err := tx.Run(cypher, map[string]interface{}{
			"id": id.String(),
		})
		if err != nil {
			return nil, err
		}

		if !result.Next() {
			return nil, nil
		}

		record := result.Record()
		u := &User{}
		u.ID, err = uuid.Parse(record.GetByIndex(0).(string))
		if err != nil {
			return nil, err
		}
		return u, nil
	})
	if err != nil {
		return nil, err
	}
	return result.(*User), nil
}
