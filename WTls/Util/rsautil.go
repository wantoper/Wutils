package Util

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
)

func GetPublicKey(certPath string) (*rsa.PublicKey, error) {
	certBytes, err := os.ReadFile(certPath)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(certBytes)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, err
	}
	return cert.PublicKey.(*rsa.PublicKey), nil
}

// 从私钥文件读取私钥
func GetPrivateKey(keyPath string) (*rsa.PrivateKey, error) {
	keyBytes, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(keyBytes)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}
	return x509.ParsePKCS1PrivateKey(block.Bytes)
}

// 加密函数
func Encrypt_RSA(plaintext []byte, publicKey *rsa.PublicKey) ([]byte, error) {
	return rsa.EncryptPKCS1v15(rand.Reader, publicKey, plaintext)
}

// 解密函数
func Decrypt_RSA(ciphertext []byte, privateKey *rsa.PrivateKey) ([]byte, error) {
	return rsa.DecryptPKCS1v15(rand.Reader, privateKey, ciphertext)
}

func PublickeyToBytes(publicKey *rsa.PublicKey) ([]byte, error) {
	return x509.MarshalPKIXPublicKey(publicKey)
}

func BytesToPublicKey(pubKeyBytes []byte) (*rsa.PublicKey, error) {
	pubKey, err := x509.ParsePKIXPublicKey(pubKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}
	if rsaPubKey, ok := pubKey.(*rsa.PublicKey); ok {
		return rsaPubKey, nil
	}

	return nil, fmt.Errorf("not an RSA public key")
}
