package main

import (
	"WUtils/Util"
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
	randonKey, _ := getRandonKey()
	//randonKey := []byte("12345678901234567890123456789012")
	fmt.Println("randonKey:", randonKey)
	key, _ := Util.GetPublicKey("server.crt")
	bytes, _ := Util.Encrypt_Rsa(randonKey, key)

	fmt.Println(bytes)

	//encryptedData, err := encrypt_aes(rsaPBK, randonKey)
	//if err != nil {
	//	panic(err)
	//}
	//
	//fmt.Println("encryptedData:", string(encryptedData))
	//
	//decryptedData, _ := decrypt_aes(encryptedData, randonKey)
	//fmt.Println("decryptedData:", string(decryptedData))

	pri_Key, _ := Util.GetPrivateKey("server.key")
	rsa, _ := Util.Decrypt_Rsa(bytes, pri_Key)
	fmt.Println(rsa)

}
