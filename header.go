package requestid

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/akm/slogctx"
)

type HeaderOptions struct {
	logAttr        string
	provider       provider
	responseSetter responseSetter
}

func newHeaderOptions() *HeaderOptions {
	name := "X-Request-ID"
	return &HeaderOptions{
		logAttr: "req_id",
		provider: selectNotEmpty(
			requestIDProvider(name),
			generatorProvider(IDGeneratorDefault),
		),
		responseSetter: newResponseSetter(name),
	}
}

type headerContextKey struct {
	name string
}

type HeaderOption = func(*HeaderOptions)

type Header struct {
	logAttr        string
	provider       provider
	responseSetter responseSetter
	ctxKey         headerContextKey
}

func NewHeader(opts ...HeaderOption) *Header {
	options := newHeaderOptions()
	for _, optFunc := range opts {
		optFunc(options)
	}
	return newHeader(options)
}

func newHeader(options *HeaderOptions) *Header {
	return &Header{
		logAttr:        options.logAttr,
		provider:       options.provider,
		responseSetter: options.responseSetter,
		ctxKey:         headerContextKey{name: options.logAttr},
	}
}

func (h *Header) Wrap(hdlr http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := h.provider(r)
		ctx := h.newContext(r.Context(), requestID)
		h.responseSetter(w, requestID)
		hdlr.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h *Header) newContext(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, h.ctxKey, requestID)
}

func (h *Header) Get(ctx context.Context) string {
	if value, ok := ctx.Value(h.ctxKey).(string); ok {
		return value
	}
	return ""
}

func (h *Header) addRecordConvTo(slogctxNS *slogctx.Namespace) {
	slogctxNS.AddRecordConv(func(ctx context.Context, rec slog.Record) slog.Record {
		requestID := h.Get(ctx)
		if requestID != "" {
			rec.Add(h.logAttr, requestID)
		}
		return rec
	})
}
