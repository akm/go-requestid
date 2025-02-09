package requestid

import "github.com/akm/slogctx"

// Options represents the options for the request ID middleware and logging request ID with slog.
type Options struct {
	logAttr          string
	generator        generator
	requestHeader    string
	responseHeader   string
	slogctxNamespace *slogctx.Namespace
}

func newDefaultOptions() *Options {
	return &Options{
		logAttr:          "req_id",
		generator:        defaultGenerator,
		requestHeader:    "",
		responseHeader:   "X-Request-ID",
		slogctxNamespace: slogctx.Default(),
	}
}

// Option represents the option for the request ID middleware and logging request ID with slog.
type Option func(o *Options)

// LogAttr sets the attribute name for logging request ID with slog.
func LogAttr(attr string) Option {
	return func(o *Options) { o.logAttr = attr }
}

// Generator sets the request ID generator.
func Generator(g generator) Option {
	return func(o *Options) { o.generator = g }
}

// RequestHeader sets the request header name for the request ID.
// Set "X-Cloud-Trace-Context" for Google Cloud https://cloud.google.com/trace/docs/trace-context?hl=ja#http-requests
// Set "X-Amzn-Trace-Id" for AWS https://docs.aws.amazon.com/ja_jp/elasticloadbalancing/latest/application/load-balancer-request-tracing.html
func RequestHeader(h string) Option {
	return func(o *Options) { o.requestHeader = h }
}

// ResponseHeader sets the response header name for the request ID.
// Skip setting the response header by setting an empty string.
func ResponseHeader(h string) Option {
	return func(o *Options) { o.responseHeader = h }
}

// SlogwNamespace sets the slogctx.Namespace for logging request ID with slog.
func SlogwNamespace(ns *slogctx.Namespace) Option {
	return func(o *Options) { o.slogctxNamespace = ns }
}
