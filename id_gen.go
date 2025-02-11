package requestid

import (
	"context"
	"log/slog"
)

// Returns a random ID generator that generates a string of length characters
// using randInt to generate random integers such as Int function from math/rand/v2 package.
func RandIntIDGenerator(
	randInt func() int,
	letters []byte,
	length int,
) func() string {
	lenLetters := len(letters)
	return func() string {
		b := make([]byte, length)
		for i := range b {
			b[i] = letters[randInt()%lenLetters]
		}
		return string(b)
	}
}

// Returns a random ID generator that generates a string of length characters
// using randRead to generate random bytes such as Read function from crypto/rand package.
func RandReadIDGenerator(
	randRead func(b []byte) (n int, err error),
	letters []byte,
	length int,
) func() (string, error) {
	lenLetters := len(letters)
	return func() (string, error) {
		b := make([]byte, length)
		if _, err := randRead(b); err != nil {
			return "", err
		}
		for i := range b {
			b[i] = letters[int(b[i])%lenLetters]
		}
		return string(b), nil
	}
}

// IDGenErrorSuppressor returns an ID generator that suppresses errors.
// If an error occurs, the recover function is called with the error and the result is returned.
func IDGenErrorSuppressor(idGen func() (string, error), recoveryFunc func(error) string) func() string {
	return func() string {
		id, err := idGen()
		if err != nil {
			return recoveryFunc(err)
		}
		return id
	}
}

// ErrorLoggingRecoveryFunc returns a recovery function that logs an error with the specified log level.
func ErrorLoggingRecoveryFunc(logLevel slog.Level, alt string) func(error) string {
	return func(err error) string {
		slog.Log(context.Background(), logLevel, "id generation error", "error", err, "alt", alt)
		return alt
	}
}
