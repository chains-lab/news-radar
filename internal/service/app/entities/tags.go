package entities

import (
	"github.com/recovery-flow/news-radar/internal/service/data"
	"github.com/sirupsen/logrus"
)

type Tags interface {
}

type tags struct {
	data *data.Tags
	log  *logrus.Entry
}
