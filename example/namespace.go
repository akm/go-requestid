package main

import (
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/akm/go-requestid"
)

func namespace() {
	helloHandler := func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		slog.InfoContext(ctx, "Start")
		defer slog.InfoContext(ctx, "End")
		io.WriteString(w, "Hello, world!\n") // nolint: errcheck
	}

	ns := requestid.New()
	ns.SlogwNamespace.AddRecordConv(requestid.SlogwPrepareFunc("requestid"))

	slog.SetDefault(ns.SlogwNamespace.New(slog.NewTextHandler(os.Stdout, nil)))

	http.Handle("/hello", ns.Wrap(http.HandlerFunc(helloHandler)))
	fmt.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
