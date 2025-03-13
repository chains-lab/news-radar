package infra

import (
	"github.com/recovery-flow/news-radar/internal/config"
	"github.com/recovery-flow/news-radar/internal/service/infra/data"
	"github.com/recovery-flow/news-radar/internal/service/infra/events/producer"
	"github.com/sirupsen/logrus"
)

type Infra struct {
	Kafka producer.Producer
	Data  *data.Data
}

func NewInfra(cfg *config.Config, log *logrus.Logger) (*Infra, error) {
	eve := producer.NewProducer(cfg)
	NewData, err := data.NewData(cfg, log)
	if err != nil {
		return nil, err
	}
	return &Infra{
		Kafka: eve,
		Data:  NewData,
	}, nil
}
