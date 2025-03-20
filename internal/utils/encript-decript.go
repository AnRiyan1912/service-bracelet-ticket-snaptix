package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base32"
)

// Generate 6-character short code from HMAC
func GenerateShortCode(data, key []byte) string {
	mac := hmac.New(sha256.New, key)
	mac.Write(data)
	hash := mac.Sum(nil)

	// Encode ke Base32 dan ambil 6 karakter pertama
	return base32.StdEncoding.EncodeToString(hash)[:6]
}

// Verify HMAC
func VerifyShortCode(data, key []byte, shortCode string) bool {
	expectedCode := GenerateShortCode(data, key)
	return expectedCode == shortCode
}
