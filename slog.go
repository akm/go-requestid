package requestid

import (
	"context"
	"log/slog"

	"github.com/akm/slogwrap"
)

func RegisterSlogHandle(key string) {
	slogwrap.Register(
		func(orig slogwrap.HandleFunc) slogwrap.HandleFunc {
			return func(ctx context.Context, rec slog.Record) error {
				requestID := Get(ctx)
				if requestID != "" {
					rec.Add(key, requestID)
				}
				return orig(ctx, rec)
			}
		},
	)
}
