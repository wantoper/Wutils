package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
)

func getRandonKey() ([]byte, error) {
	aesKey := make([]byte, 32) // AES-256
	_, err := io.ReadFull(rand.Reader, aesKey)
	if err != nil {
		return nil, err
	}
	return aesKey, nil
}

func encrypt_aes(data, key []byte) ([]byte, error) {
	block, _ := aes.NewCipher(key)
	gcm, _ := cipher.NewGCM(block)
	nonce := make([]byte, gcm.NonceSize())
	io.ReadFull(rand.Reader, nonce)
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext, nil
}

func decrypt_aes(ciphertext, key []byte) ([]byte, error) {
	block, _ := aes.NewCipher(key)
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

func main() {
	randonKey, err := getRandonKey()
	if err != nil {
		panic(err)
	}

	key, _ := getPublicKey("server.crt")
	bytes, err := encrypt(randonKey, key)

	fmt.Println(string(bytes))

	//encryptedData, err := encrypt_aes(rsaPBK, randonKey)
	//if err != nil {
	//	panic(err)
	//}
	//
	//fmt.Println("encryptedData:", string(encryptedData))
	//
	//decryptedData, _ := decrypt_aes(encryptedData, randonKey)
	//fmt.Println("decryptedData:", string(decryptedData))

}
