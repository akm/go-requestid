package requestid

func Default() *Factory {
	return &Factory{
		generator:      defaultGenerator,
		requestHeader:  "X-Request-ID",
		responseHeader: "X-Request-ID",
	}
}
