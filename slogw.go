package requestid

import (
	"context"
	"log/slog"
)

func SlogwPrepareFunc(key string) func(ctx context.Context, rec slog.Record) slog.Record {
	return func(ctx context.Context, rec slog.Record) slog.Record {
		requestID := Get(ctx)
		if requestID != "" {
			rec.Add(key, requestID)
		}
		return rec
	}
}
