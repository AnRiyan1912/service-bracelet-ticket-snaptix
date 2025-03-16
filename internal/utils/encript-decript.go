package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

func EncryptAESGCM(plainText, key string) (string, error) {
	keyBytes := make([]byte, 32)
	copy(keyBytes, key)

	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	cipherText := aesGCM.Seal(nil, nonce, []byte(plainText), nil)
	finalCipher := append(nonce, cipherText...)
	return base64.StdEncoding.EncodeToString(finalCipher), nil
}

func DecryptAESGCM(encryptedBase64, key string) (string, error) {
	// Decode Base64
	encryptedData, err := base64.StdEncoding.DecodeString(encryptedBase64)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64: %v", err)
	}

	// Pastikan key memiliki panjang 32 byte (AES-256)
	keyBytes := make([]byte, 32)
	copy(keyBytes, key)

	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %v", err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %v", err)
	}

	// Pisahkan nonce (12 byte pertama)
	nonceSize := aesGCM.NonceSize()
	if len(encryptedData) < nonceSize {
		return "", fmt.Errorf("invalid encrypted data")
	}

	nonce := encryptedData[:nonceSize]
	cipherText := encryptedData[nonceSize:]

	// Dekripsi data
	plainText, err := aesGCM.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt: %v", err)
	}

	return string(plainText), nil
}
