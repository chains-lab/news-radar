package entities

import (
	"github.com/recovery-flow/news-radar/internal/service/data"
	"github.com/sirupsen/logrus"
)

type User interface {
}

type user struct {
	data *data.Users
	log  *logrus.Entry
}
