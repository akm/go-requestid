package requestid_test

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"

	"github.com/akm/go-requestid"
	"github.com/akm/slogctx"
)

func Example() {
	helloHandler := func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		slog.InfoContext(ctx, "Start")
		defer slog.InfoContext(ctx, "End")
		io.WriteString(w, "Hello, world!\n") //nolint:errcheck
	}

	// This is a demo. Call slogctx.Register from init function in your application.
	slogctx.Add(requestid.SlogwPrepareFunc("requestid"))
	slog.SetDefault(slogctx.New(slog.NewTextHandler(os.Stdout, testOptions)))

	ts := httptest.NewServer(
		requestid.Wrap(http.HandlerFunc(helloHandler)),
	)
	defer ts.Close()

	resp, _ := http.Get(ts.URL) //nolint:errcheck
	defer resp.Body.Close()     //nolint:govet

	b, _ := io.ReadAll(resp.Body) //nolint:errcheck
	fmt.Println(string(b))

	// Output logs with requestid with generated ID
	// level=INFO msg=Start requestid=k3pGQp5T
	// level=INFO msg=End requestid=k3pGQp5T
	// Hello, world!
}
