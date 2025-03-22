package entities

import (
	"github.com/recovery-flow/news-radar/internal/service/data"
	"github.com/sirupsen/logrus"
)

type Themes interface {
}

type theme struct {
	data *data.Themes
	log  *logrus.Entry
}
