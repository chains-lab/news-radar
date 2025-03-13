package service

import (
	"github.com/recovery-flow/news-radar/internal/config"
	"github.com/recovery-flow/news-radar/internal/service/domain"
	"github.com/recovery-flow/news-radar/internal/service/infra"
	"github.com/sirupsen/logrus"
)

type Service struct {
	Config *config.Config
	Domain domain.Domain
	Log    *logrus.Logger
}

func NewService(cfg *config.Config, log *logrus.Logger) (*Service, error) {
	inf, err := infra.NewInfra(cfg, log)
	if err != nil {
		return nil, err
	}
	dmn, err := domain.NewDomain(inf, log)
	if err != nil {
		return nil, err
	}

	return &Service{
		Config: cfg,
		Log:    log,
		Domain: dmn,
	}, nil
}
