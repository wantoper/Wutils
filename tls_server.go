package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
)

func main() {
	cert, err := tls.LoadX509KeyPair("server.crt", "server.key")
	if err != nil {
		log.Fatalf("服务器证书加载失败: %v", err)
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	listener, err := tls.Listen("tcp", ":4443", config)
	if err != nil {
		log.Fatalf("启动 TLS 服务器失败: %v", err)
	}
	defer listener.Close()

	fmt.Println("TLS 服务器已启动，监听端口 :4443")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("接受连接失败: %v", err)
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			log.Printf("读取数据失败: %v", err)
			return
		}

		fmt.Printf("收到消息: %s\n", string(buffer[:n]))

		// 发送回复
		response := []byte("服务器已收到消息：" + string(buffer[:n]))
		if _, err := conn.Write(response); err != nil {
			log.Printf("发送响应失败: %v", err)
			return
		}
	}
}
