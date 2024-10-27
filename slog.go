package requestid

import (
	"context"
	"log/slog"

	"github.com/akm/slogw"
)

func RegisterSlogHandle(key string) {
	slogw.Register(
		func(orig slogw.HandleFunc) slogw.HandleFunc {
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
