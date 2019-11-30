package rand

import (
	"fmt"
	"math/rand"
	"time"
)

const source_charset = "0123456789" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ"

// Make sure it is pseudo-random by using a new seed on startup
var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

func stringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func randomString(length int) string {
	return stringWithCharset(length, source_charset)
}

func GenerateNewSymbol(original string) string {
	return fmt.Sprintf("%s-%s", original, randomString(3))
}
