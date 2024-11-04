package shortener

import (
	"crypto/rand"
	"encoding/base64"
)

func GeneratorShortURL() string {
	cap_len := 6
	b := make([]byte, cap_len)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return base64.URLEncoding.EncodeToString(b)
}
