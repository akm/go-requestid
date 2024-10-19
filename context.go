package requestid

import "context"

type ctxKeyType struct{}

var ctxKey = ctxKeyType{}

func set(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, ctxKey, requestID)
}

func Get(ctx context.Context) string {
	if value, ok := ctx.Value(ctxKey).(string); ok {
		return value
	}
	return ""
}
