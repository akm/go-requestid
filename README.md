# go-requestid

![CI](https://github.com/akm/go-requestid/actions/workflows/ci.yml/badge.svg)

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

### with slog

```golang
import "github.com/akm/slogw"
```

Setup logger

```golang
    slog.SetDefault(slogw.New(slog.NewTextHandler(os.Stdout, nil)))
```

you can use slog.NewJSONHandler instead of slog.NewTextHandler.

And setup slog Handler for requestid.

```golang
func init() {
	requestid.RegisterSlogHandle("requestid")
}
```

Then the server log includes requestid.

```shell
$ go run ./example
Server started at :8080
time=2024-10-19T23:29:52.650+09:00 level=INFO msg=Start requestid=eVxKgnfE
time=2024-10-19T23:29:52.651+09:00 level=INFO msg=End requestid=eVxKgnfE
time=2024-10-19T23:29:53.389+09:00 level=INFO msg=Start requestid=MJsEk1JG
time=2024-10-19T23:29:53.389+09:00 level=INFO msg=End requestid=MJsEk1JG
```

See [example/main.go](./example/main.go) for more details.
