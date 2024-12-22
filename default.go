package requestid

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/akm/slogw"
)

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
