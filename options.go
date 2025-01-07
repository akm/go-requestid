package requestid

import "github.com/akm/slogw"

type Options struct {
	Generator      generator
	RequestHeader  string
	ResponseHeader string
	SlogwNamespace *slogw.Namespace
}

func newDefaultOptions() *Options {
	return &Options{
		Generator:      defaultGenerator,
		RequestHeader:  "",
		ResponseHeader: "X-Request-ID",
		SlogwNamespace: slogw.Default(),
	}
}

type Option func(o *Options)

func Generator(g generator) Option {
	return func(o *Options) {
		o.Generator = g
	}
}

// Set X-Cloud-Trace-Context for Google Cloud https://cloud.google.com/trace/docs/trace-context?hl=ja#http-requests
// Set X-Amzn-Trace-Id for AWS https://docs.aws.amazon.com/ja_jp/elasticloadbalancing/latest/application/load-balancer-request-tracing.html
func RequestHeader(h string) Option {
	return func(o *Options) {
		o.RequestHeader = h
	}
}

func ResponseHeader(h string) Option {
	return func(o *Options) {
		o.ResponseHeader = h
	}
}

func SlogwNamespace(ns *slogw.Namespace) Option {
	return func(o *Options) {
		o.SlogwNamespace = ns
	}
}
