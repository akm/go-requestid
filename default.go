package requestid

func Default() *Options {
	return &Options{
		Generator:      defaultGenerator,
		RequestHeader:  "X-Request-ID",
		ResponseHeader: "X-Request-ID",
	}
}
