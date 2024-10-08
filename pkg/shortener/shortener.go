package shortener

import (
	"crypto/rand"
	"encoding/base64"
)

func GeneratorShortURL() string {
	b := make([]byte, 6)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}
