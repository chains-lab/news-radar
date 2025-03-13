package domain

import (
	"github.com/recovery-flow/news-radar/internal/service/infra"
	"github.com/sirupsen/logrus"
)

type Domain interface {
}

type domain struct {
	Infra *infra.Infra
	log   *logrus.Logger
}

func NewDomain(inf *infra.Infra, log *logrus.Logger) (Domain, error) {
	return &domain{
		Infra: inf,
		log:   log,
	}, nil
}
