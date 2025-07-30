package main

import (
	"WUtils/WTLS"
	"log"
)

func main() {
	client, _ := WTLS.NewWTLSClient()
	err := client.Connect(":8080")
	if err != nil {
		log.Fatal(err)
	}
	data := "Hello, WTLS Server!"
	_, err = client.Write([]byte(data))
	if err != nil {
		log.Fatal("发送数据失败:", err)
	} else {
		log.Println("发送数据成功:", data)
	}
}
