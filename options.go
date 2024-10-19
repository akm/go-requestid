package requestid

type Options struct {
	Generator      generator
	RequestHeader  string
	ResponseHeader string
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
