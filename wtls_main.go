package main

import (
	"WUtils/WTLS"
	"bufio"
	"fmt"
)

func main() {
	server, err := WTLS.NewWTLSServer()
	if err != nil {
		panic("创建 WTLS 服务器失败: " + err.Error())
	}
	err = server.Listen(":4433")
	if err != nil {
		panic("监听失败: " + err.Error())
	}

	for {
		fmt.Println("等待客户端连接...")
		conn, err := server.Accept()
		fmt.Println("新连接来自")
		if err != nil {
			panic("接受连接失败: " + err.Error())
		}

		reader := bufio.NewReader(conn)
		for {
			// 读取客户端发送的数据，以换行符为结束标志
			data, err := reader.ReadString('\n')
			if err != nil {
				fmt.Printf("读取数据失败: %v\n", err)
				break
			}

			// 打印接收到的数据
			fmt.Printf("接收到: %s", data)

			// 可选：向客户端发送响应
			response := "已收到数据: " + data
			_, err = conn.Write([]byte(response))
			if err != nil {
				fmt.Printf("向发送响应失败: %v\n", err)
				break
			}
		}
	}

}
