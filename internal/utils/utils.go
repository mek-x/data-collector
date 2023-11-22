package utils

import (
	"math/rand"
	"os"
	"regexp"
	"strings"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func stringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func RandomString(length int) string {
	return stringWithCharset(length, charset)
}

func ReplaceWithEnvVars(in string) string {
	re := regexp.MustCompile(`%%[a-zA-Z][0-9a-zA-Z_]+%%`)
	out := re.ReplaceAllStringFunc(in, func(s string) string {
		s = os.Getenv(strings.Trim(s, "%"))
		return s
	})
	return out
}
