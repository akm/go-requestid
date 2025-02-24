package requestid_test

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"

	"github.com/akm/go-requestid"
)

func ExampleMiddleware() {
	initializeSlogctx()

	mw := requestid.New(
		requestid.LogAttr("req-id"),
		requestid.Generator(func() string { return "(id-for-test)" }),
	)

	logger := mw.NewLogger(slog.NewTextHandler(os.Stdout, testOptions))

	helloHandler := func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		logger.InfoContext(ctx, "Start")
		defer logger.InfoContext(ctx, "End")
		io.WriteString(w, "Hello, world!\n") //nolint:errcheck
	}

	ts := httptest.NewServer(
		mw.Wrap(http.HandlerFunc(helloHandler)),
	)
	defer ts.Close()

	resp, _ := http.Get(ts.URL) //nolint:errcheck
	defer resp.Body.Close()     //nolint:govet

	b, _ := io.ReadAll(resp.Body) //nolint:errcheck
	fmt.Println(string(b))

	// Output:
	// level=INFO msg=Start req-id=(id-for-test)
	// level=INFO msg=End req-id=(id-for-test)
	// Hello, world!
}
