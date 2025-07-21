package sender

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

// Функция вычисления хэша
func calculateHash(data, key []byte) string {
	h := hmac.New(sha256.New, key)
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}
