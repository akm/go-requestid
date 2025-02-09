package requestid

import "github.com/akm/slogctx"

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

type Option func(o *Options)

func LogAttr(attr string) Option {
	return func(o *Options) {
		o.logAttr = attr
	}
}

func Generator(g generator) Option {
	return func(o *Options) {
		o.generator = g
	}
}

// Set X-Cloud-Trace-Context for Google Cloud https://cloud.google.com/trace/docs/trace-context?hl=ja#http-requests
// Set X-Amzn-Trace-Id for AWS https://docs.aws.amazon.com/ja_jp/elasticloadbalancing/latest/application/load-balancer-request-tracing.html
func RequestHeader(h string) Option {
	return func(o *Options) {
		o.requestHeader = h
	}
}

func ResponseHeader(h string) Option {
	return func(o *Options) {
		o.responseHeader = h
	}
}

func SlogwNamespace(ns *slogctx.Namespace) Option {
	return func(o *Options) {
		o.slogctxNamespace = ns
	}
}
