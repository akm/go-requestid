package requestid

func Default() *Options {
	return &Options{
		generator:      defaultGenerator,
		requestHeader:  "X-Request-ID",
		responseHeader: "X-Request-ID",
	}
}
