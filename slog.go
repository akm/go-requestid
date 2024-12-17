package requestid

import (
	"context"
	"log/slog"

	"github.com/akm/slogw"
)

func RegisterSlogHandle(key string) {
	slogw.Register(
		func(ctx context.Context, rec slog.Record) slog.Record {
			requestID := Get(ctx)
			if requestID != "" {
				rec.Add(key, requestID)
			}
			return rec
		},
	)
}
