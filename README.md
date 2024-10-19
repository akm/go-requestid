# go-requestid

## Install

```
go get github.com/akm/go-requestid
```

## How to use

### import

```
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
