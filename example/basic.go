package main

import (
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/akm/go-requestid"
	"github.com/akm/slogctx"
)

func basic() {
	// This is a demo. Call slogctx.Register from init function in your application.
	slogctx.Add(requestid.Handle("requestid"))

	helloHandler := func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		slog.InfoContext(ctx, "Start")
		defer slog.InfoContext(ctx, "End")
		io.WriteString(w, "Hello, world!\n") // nolint: errcheck
	}

	slog.SetDefault(slogctx.New(slog.NewTextHandler(os.Stdout, nil)))

	http.Handle("/hello", requestid.Wrap(http.HandlerFunc(helloHandler)))
	fmt.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
