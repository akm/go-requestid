# go-requestid

[![CI](https://github.com/akm/go-requestid/actions/workflows/ci.yml/badge.svg)](https://github.com/akm/go-requestid/actions/workflows/ci.yml)
[![codecov](https://codecov.io/github/akm/go-requestid/graph/badge.svg?token=9BcanbSLut)](https://codecov.io/github/akm/go-requestid)
[![Enabled Linters](https://img.shields.io/badge/dynamic/yaml?url=https%3A%2F%2Fraw.githubusercontent.com%2Fakm%2Fgo-requestid%2Frefs%2Fheads%2Fmain%2F.project.yaml&query=%24.linters&label=enabled%20linters&color=%2317AFC2)](.golangci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/akm/go-requestid)](https://goreportcard.com/report/github.com/akm/go-requestid)
[![Documentation](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/akm/go-requestid)
[![license](https://img.shields.io/github/license/akm/go-requestid)](./LICENSE)

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

### with slog Logger

Setup logger

```golang
    logger := requestid.NewLoger(slog.NewTextHandler(os.Stdout, nil))
```

And setup slog Handler for requestid.

See [example_test.go](./example_test.go) for more details.
