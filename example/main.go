package main

import (
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/akm/go-requestid"
	"github.com/akm/slogwrap"
)

func main() {
	helloHandler := func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		slog.InfoContext(ctx, "Start")
		defer slog.InfoContext(ctx, "End")
		io.WriteString(w, "Hello, world!\n") // nolint: errcheck
	}

	slog.SetDefault(slogwrap.New(slog.NewTextHandler(os.Stdout, nil)))

	http.Handle("/hello", requestid.Wrap(http.HandlerFunc(helloHandler)))
	fmt.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func init() {
	requestid.RegisterSlogHandle("requestid")
}
