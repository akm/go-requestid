package requestid

import (
	"crypto/rand"
	"encoding/base64"
)

type generator = func() string

// RandomGenerator returns a generator that generates a random string of the specified size.
func RandomGenerator(size int) generator {
	return func() string {
		b := make([]byte, size)
		rand.Read(b) //nolint:errcheck
		return base64.RawURLEncoding.EncodeToString(b)
	}
}

var defaultGenerator = RandomGenerator(6)
