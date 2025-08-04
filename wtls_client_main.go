package main

import (
	"WUtils/WTls"
	"fmt"
)

func main() {
	client, err := WTls.Dial(":4443")
	if err != nil {
		panic(err)
	}
	write, err := client.Write([]byte("测试发送的消息测试发送的消息测试发送的消息测试发送的消息测试发送的消息测试发送的消息测试发送的消息测试发送的消息测试发送的消息测试发送的消息测试发送的消息测试发送的消息测试发送的消息测试发送的消息测试发送的消息测试发送的消息测试发送的消息测试发送的消息测试发送的消息测试发送的消息测试发送的消息测试发送的消息测试发送的消息测试发送的消息测试发送的消息测试发送的消息测试发送的消息测试发送的消息"))
	if err != nil {
		fmt.Printf("Write error: %v\n", err)
		return
	}
	fmt.Printf("Write %d bytes\n", write)
	client.Close()
}
