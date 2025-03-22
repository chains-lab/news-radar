package entities

import (
	"github.com/recovery-flow/news-radar/internal/service/data"
	"github.com/sirupsen/logrus"
)

type Articles interface {
}

type articles struct {
	data data.Article
	log  *logrus.Entry
}
