package token_generator

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func Generate() string {
	b := make([]byte, 32)
	rand.Read(b)

	return base64.URLEncoding.EncodeToString([]byte(fmt.Sprintf("%x", b)))
}
