package requestid_test

import (
	"log/slog"

	"github.com/akm/slogctx"
)

var testOptions = &slog.HandlerOptions{
	ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
		if a.Key == slog.TimeKey {
			return slog.Attr{}
		}
		return a
	},
}

func initializeSlogctx() {
	slogctx.SetDefault(slogctx.NewNamespace())
}
