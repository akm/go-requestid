package requestid

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/akm/slogw"
)

var degaultNamespace = newFactory(Default())

func DefaultNamespace() *Namespace {
	return degaultNamespace
}

func SetDefaultNamespace(ns *Namespace) {
	degaultNamespace = ns
}

func RegisterSlogHandle(key string) {
	slogw.Register(func(ctx context.Context, rec slog.Record) slog.Record {
		requestID := Get(ctx)
		if requestID != "" {
			rec.Add(key, requestID)
		}
		return rec
	})
}

func Wrap(next http.Handler, opts ...Option) http.Handler {
	options := Default()
	for _, optFunc := range opts {
		optFunc(options)
	}
	return newFactory(options).Wrap(next)
}
