package utils

import (
	"crypto/rand"
	"encoding/base64"
)

func GenerateRandomState() string {
	b := make([]byte, 16)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}
