package requestid

import (
	"crypto/rand"
	"encoding/base64"
)

type Generator = func() string

func RandomGenerator(size int) Generator {
	return func() string {
		b := make([]byte, size)
		rand.Read(b) //nolint:errcheck
		return base64.RawURLEncoding.EncodeToString(b)
	}
}

var defaultGenerator = RandomGenerator(6)
