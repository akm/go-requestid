# go-requestid

[![CI](https://github.com/akm/go-requestid/actions/workflows/ci.yml/badge.svg)](https://github.com/akm/go-requestid/actions/workflows/ci.yml)
[![codecov](https://codecov.io/github/akm/go-requestid/graph/badge.svg?token=9BcanbSLut)](https://codecov.io/github/akm/go-requestid)
[![Enabled Linters](https://img.shields.io/badge/dynamic/yaml?url=https%3A%2F%2Fraw.githubusercontent.com%2Fakm%2Fgo-requestid%2Frefs%2Fheads%2Fmain%2F.project.yaml&query=%24.linters&label=enabled%20linters&color=%2317AFC2)](.golangci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/akm/go-requestid)](https://goreportcard.com/report/github.com/akm/go-requestid)
[![Documentation](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/akm/go-requestid)
[![license](https://img.shields.io/github/license/akm/go-requestid)](./LICENSE)

## Overview

`github.com/akm/go-requestid` is a go module for request ID. The package name is `requestid`.
It works in the following sequence for each request processing:

1. requestid gets request ID from HTTP request header or generates request ID
2. requestid create a new [context.Context](https://pkg.go.dev/context#Context) with request ID
3. When your application logs with [slog.Logger](https://pkg.go.dev/log/slog#Logger) by requestid, the request ID is added to a log record as an attribute.

## Install

```shell
go get github.com/akm/go-requestid
```

## How to use

### import

```golang
import "github.com/akm/go-requestid"
```

### Easy way

```golang
handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    // Get the request ID by using requestid.Get in your http handler
    requestID := requestid.Get(r.Context())
    w.Write([]byte(fmt.Sprintf("Request ID", requestID)))
    w.WriteHeader(http.StatusOK)
})

// Wrap your http handler in order to process requestid
handlerWithRequestID := requestid.Wrap(handler)
```

### With slog Logger

Setup logger

```golang
logger := requestid.NewLoger(slog.NewTextHandler(os.Stdout, nil))
```

And setup slog Handler for requestid.

See [example_test.go](./example_test.go) for more details.

## X-Request-ID

`X-Request-ID` is an unofficial HTTP request/response header. But it is supported by some services, middlewares, frameworks and libraries.

- [http.dev / X-Request-ID](https://http.dev/x-request-id)
- [Envoy / Tracing](https://www.envoyproxy.io/docs/envoy/latest/intro/arch_overview/observability/tracing)
- [nginx / $request_id](https://nginx.org/en/docs/http/ngx_http_core_module.html#var_request_id)
- [Heroku / HTTP Request IDs](https://devcenter.heroku.com/articles/http-request-id)
- [Ruby on Rails / Action Dispatch RequestId](https://api.rubyonrails.org/classes/ActionDispatch/RequestId.html)
- [django-log-request-id](https://github.com/dabapps/django-log-request-id)

How to use `X-Request-ID` varies depending on the situation.
`X-Request-ID` header can be trusted if the application runs behind a proxy such as Envoy or nginx that generates that header.
`X-Request-ID` header is unreliable if your application communicates with the client directly or if your proxy does not modify that header.
In the latter case, you should consider using `X-Client-Request-ID`.

### ID generators

`requestid` provides two ID generators which work with the following packages:

- math/rand/v2
- crypto/rand
