package requestid

import (
	"context"
	"log/slog"

	"github.com/akm/slogctx"
)

func SlogwPrepareFunc(key string) slogctx.RecordConv {
	return func(ctx context.Context, rec slog.Record) slog.Record {
		requestID := Get(ctx)
		if requestID != "" {
			rec.Add(key, requestID)
		}
		return rec
	}
}
