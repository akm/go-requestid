package requestid

import (
	"crypto/rand"
	"log/slog"
)

type generator = func() string

const defaultIDLength = 16

var defaultIDLetters = []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_")

// IDGeneratorDefault is the default ID generator.
var IDGeneratorDefault = IDGenErrorSuppressor(
	RandReadIDGenerator(rand.Read, defaultIDLetters, defaultIDLength),
	ErrorLoggingRecoveryFunc(slog.LevelWarn, "id-gen-error"),
)
