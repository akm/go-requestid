package requestid

import (
	"context"
	"log/slog"

	"github.com/akm/slogw"
)

func SlogwPrepareFunc(key string) slogw.HandlePrepareFunc {
	return func(ctx context.Context, rec slog.Record) slog.Record {
		requestID := Get(ctx)
		if requestID != "" {
			rec.Add(key, requestID)
		}
		return rec
	}
}
