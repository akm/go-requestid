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

func Example() {
	logger := requestid.NewLogger(slog.NewTextHandler(os.Stdout, testOptions))

	helloHandler := func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		logger.InfoContext(ctx, "Start")
		defer logger.InfoContext(ctx, "End")
		io.WriteString(w, "Hello, world!\n") //nolint:errcheck
	}

	ts := httptest.NewServer(
		requestid.Wrap(http.HandlerFunc(helloHandler)),
	)
	defer ts.Close()

	resp, _ := http.Get(ts.URL) //nolint:errcheck
	defer resp.Body.Close()     //nolint:govet

	b, _ := io.ReadAll(resp.Body) //nolint:errcheck
	fmt.Println(string(b))

	// Output logs with req_id with generated ID
	// level=INFO msg=Start req_id=k3pGQp5T
	// level=INFO msg=End req_id=k3pGQp5T
	// Hello, world!
}
