package requestid

func Default() *Options {
	return &Options{
		Generator: defaultGenerator,
		// Set X-Cloud-Trace-Context for Google Cloud https://cloud.google.com/trace/docs/trace-context?hl=ja#http-requests
		// Set X-Amzn-Trace-Id for AWS https://docs.aws.amazon.com/ja_jp/elasticloadbalancing/latest/application/load-balancer-request-tracing.html
		RequestHeader:  "",
		ResponseHeader: "X-Request-ID",
	}
}
