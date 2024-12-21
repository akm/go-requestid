package requestid

type Options struct {
	Generator      generator
	RequestHeader  string
	ResponseHeader string
}

func Default() *Options {
	return &Options{
		Generator: defaultGenerator,
		// Set X-Cloud-Trace-Context for Google Cloud https://cloud.google.com/trace/docs/trace-context?hl=ja#http-requests
		// Set X-Amzn-Trace-Id for AWS https://docs.aws.amazon.com/ja_jp/elasticloadbalancing/latest/application/load-balancer-request-tracing.html
		RequestHeader:  "",
		ResponseHeader: "X-Request-ID",
	}
}

type Option func(o *Options)

func Generator(g generator) Option {
	return func(o *Options) {
		o.Generator = g
	}
}

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
