package handlers

import (
	"context"
	"net/http"

	"github.com/recovery-flow/news-radar/internal/config"
	"github.com/recovery-flow/news-radar/internal/service/domain"
	"github.com/sirupsen/logrus"
)

type ctxKey int

const (
	logCtxKey ctxKey = iota
	domainCtxKey
	configCtxKey
)

func CtxLog(entry *logrus.Logger) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, logCtxKey, entry)
	}
}

func Log(r *http.Request) *logrus.Logger {
	return r.Context().Value(logCtxKey).(*logrus.Logger)
}

func CtxDomain(entry domain.Domain) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, domainCtxKey, entry)
	}
}

func Domain(r *http.Request) domain.Domain {
	return r.Context().Value(domainCtxKey).(domain.Domain)
}

func CtxConfig(entry *config.Config) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, configCtxKey, entry)
	}
}

func Config(r *http.Request) *config.Config {
	return r.Context().Value(configCtxKey).(*config.Config)
}
