package Util

import (
	"WUtils/WTls/consts"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
)

func GetRandonKey(version uint8) ([]byte, error) {
	var length int
	switch version {
	case consts.AES_GCM_128:
		length = 16 // AES-128
	case consts.AES_GCM_256:
		length = 32 // AES-256
	}

	aesKey := make([]byte, length)
	_, err := io.ReadFull(rand.Reader, aesKey)
	if err != nil {
		return nil, err
	}
	return aesKey, nil
}

func Encrypt_AES(data, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create AES cipher: %w", err)
	}
	gcm, _ := cipher.NewGCM(block)
	nonce := make([]byte, gcm.NonceSize())
	io.ReadFull(rand.Reader, nonce)
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext, nil
}

func Decrypt_AES(ciphertext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create AES cipher: %w", err)
	}
	gcm, _ := cipher.NewGCM(block)
	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}
