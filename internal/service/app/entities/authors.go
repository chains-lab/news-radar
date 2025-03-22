package entities

import (
	"github.com/recovery-flow/news-radar/internal/service/data"
	"github.com/sirupsen/logrus"
)

type Authors interface {
}

type authors struct {
	data *data.Authors
	log  *logrus.Logger
}
