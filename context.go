package requestid

import "context"

type ctxKeyType struct{}

var ctxKey = ctxKeyType{}

func newContext(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, ctxKey, requestID)
}

// Get returns the request ID from the context set by the middleware.
func Get(ctx context.Context) string {
	if value, ok := ctx.Value(ctxKey).(string); ok {
		return value
	}
	return ""
}
