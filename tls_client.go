package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
)

func main() {
	// 加载服务器的证书
	cert, err := ioutil.ReadFile("server.crt")
	if err != nil {
		log.Fatalf("无法读取证书文件: %v", err)
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(cert) {
		log.Fatal("无法解析证书")
	}

	config := &tls.Config{
		RootCAs:            certPool,
		InsecureSkipVerify: true, // 在开发环境中可以使用，生产环境建议设置为 false
	}

	// 连接到服务器
	conn, err := tls.Dial("tcp", "localhost:4443", config)
	if err != nil {
		log.Fatalf("连接服务器失败: %v", err)
	}
	defer conn.Close()

	fmt.Println("已成功连接到服务器")

	// 发送消息
	message := "Hello, TLS Server!"
	_, err = conn.Write([]byte(message))
	if err != nil {
		log.Fatalf("发送消息失败: %v", err)
	}

	// 接收服务器响应
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		log.Fatalf("读取响应失败: %v", err)
	}

	fmt.Printf("服务器响应: %s\n", string(buffer[:n]))
}
