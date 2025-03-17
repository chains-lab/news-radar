package infra

import (
	"github.com/recovery-flow/news-radar/internal/config"
	"github.com/recovery-flow/news-radar/internal/service/infra/data/mongodb"
	"github.com/recovery-flow/news-radar/internal/service/infra/data/neo"
	"github.com/recovery-flow/news-radar/internal/service/infra/events/producer"
	"github.com/sirupsen/logrus"
)

type Infra struct {
	Kafka producer.Producer
	Neo   Neo
	Mongo Mongo
}

type Neo struct {
	Tags     neo.Tags
	Articles neo.Articles
	Authors  neo.Authors
	Themes   neo.Themes
}

type Mongo struct {
	Articles mongodb.Articles
	Authors  mongodb.Authors
}

func NewInfra(cfg *config.Config, log *logrus.Logger) (*Infra, error) {
	eve := producer.NewProducer(cfg)
	tagNeo, err := neo.NewTags(cfg.Database.Neo4j.URI, cfg.Database.Neo4j.User, cfg.Database.Neo4j.Password)
	if err != nil {
		return nil, err
	}
	articleNeo, err := neo.NewArticles(cfg.Database.Neo4j.URI, cfg.Database.Neo4j.User, cfg.Database.Neo4j.Password)
	if err != nil {
		return nil, err
	}
	authorNeo, err := neo.NewAuthors(cfg.Database.Neo4j.URI, cfg.Database.Neo4j.User, cfg.Database.Neo4j.Password)
	if err != nil {
		return nil, err
	}
	themeNeo, err := neo.NewThemes(cfg.Database.Neo4j.URI, cfg.Database.Neo4j.User, cfg.Database.Neo4j.Password)
	if err != nil {
		return nil, err
	}

	articleMongo, err := mongodb.NewArticles(cfg.Database.Mongo.URI, cfg.Database.Mongo.Name)
	if err != nil {
		return nil, err
	}
	authorMongo, err := mongodb.NewAuthors(cfg.Database.Mongo.URI, cfg.Database.Mongo.Name)
	if err != nil {
		return nil, err
	}

	return &Infra{
		Kafka: eve,
		Neo: Neo{
			Tags:     tagNeo,
			Articles: articleNeo,
			Authors:  authorNeo,
			Themes:   themeNeo,
		},
		Mongo: Mongo{
			Articles: articleMongo,
			Authors:  authorMongo,
		},
	}, nil
}
