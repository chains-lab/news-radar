package data

import (
	"github.com/recovery-flow/news-radar/internal/config"
	"github.com/sirupsen/logrus"
)

type Data struct {
}

func NewData(cfg *config.Config, log *logrus.Logger) (*Data, error) {
	return &Data{}, nil
}
