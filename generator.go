package requestid

import (
	"crypto/rand"
	"encoding/base64"
)

type generator = func() string

func RandomGenerator(size int) generator {
	return func() string {
		b := make([]byte, size)
		rand.Read(b) //nolint:errcheck
		return base64.RawURLEncoding.EncodeToString(b)
	}
}

var defaultGenerator = RandomGenerator(6)
