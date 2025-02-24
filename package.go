package requestid

import (
	"log/slog"
	"net/http"

	"github.com/akm/slogctx"
)

var defaultNamespace = newMiddleware(newDefaultOptions())

// Default returns the default Namespace.
func Default() *Middleware {
	return defaultNamespace
}

// SetDefault sets the default Namespace.
func SetDefault(ns *Middleware) {
	defaultNamespace = ns
}

// ResetDefault resets the default Namespace to the initial state.
func ResetDefault() {
	slogctx.SetDefault(slogctx.NewNamespace())
	defaultNamespace = newMiddleware(newDefaultOptions())
}

// Wrap wraps the given http.Handler with the default Namespace.
func Wrap(next http.Handler) http.Handler {
	return defaultNamespace.Wrap(next)
}

// NewLogger returns a new logger with the default Namespace.
func NewLogger(h slog.Handler) *slog.Logger {
	return defaultNamespace.NewLogger(h)
}

// WrapSlogHandler wraps the given slog.Handler with the default Namespace.
func WrapSlogHandler(h slog.Handler) slog.Handler {
	return defaultNamespace.WrapSlogHandler(h)
}
