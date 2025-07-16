package sender

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"testing"
)

func TestCalculateHash_KnownValue(t *testing.T) {
	data := []byte("hello world")
	key := []byte("secret")
	expected := hmacSHA256Hex(data, key)
	got := calculateHash(data, key)
	if got != expected {
		t.Errorf("expected %s, got %s", expected, got)
	}
}

func TestCalculateHash_Empty(t *testing.T) {
	data := []byte("")
	key := []byte("")
	expected := hmacSHA256Hex(data, key)
	got := calculateHash(data, key)
	if got != expected {
		t.Errorf("expected %s, got %s", expected, got)
	}
}

func hmacSHA256Hex(data, key []byte) string {
	h := hmac.New(sha256.New, key)
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}
