package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
)

// 从证书文件读取公钥
func getPublicKey(certPath string) (*rsa.PublicKey, error) {
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
func getPrivateKey(keyPath string) (*rsa.PrivateKey, error) {
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
func encrypt(plaintext []byte, publicKey *rsa.PublicKey) ([]byte, error) {
	return rsa.EncryptPKCS1v15(rand.Reader, publicKey, plaintext)
}

// 解密函数
func decrypt(ciphertext []byte, privateKey *rsa.PrivateKey) ([]byte, error) {
	return rsa.DecryptPKCS1v15(rand.Reader, privateKey, ciphertext)
}

func main() {
	// 读取公钥和私钥
	publicKey, err := getPublicKey("server.crt")
	if err != nil {
		fmt.Printf("获取公钥失败: %v\n", err)
		return
	}

	privateKey, err := getPrivateKey("server.key")
	if err != nil {
		fmt.Printf("获取私钥失败: %v\n", err)
		return
	}

	// 测试加密和解密
	message := []byte("Hello, RSA加密测试!")
	fmt.Printf("原始信息: %s\n", message)

	// 加密
	encrypted, err := encrypt(message, publicKey)
	if err != nil {
		fmt.Printf("加密失败: %v\n", err)
		return
	}
	fmt.Printf("加密后的数据: %x\n", encrypted)

	// 解密
	decrypted, err := decrypt(encrypted, privateKey)
	if err != nil {
		fmt.Printf("解密失败: %v\n", err)
		return
	}
	fmt.Printf("解密后的信息: %s\n", decrypted)
}
