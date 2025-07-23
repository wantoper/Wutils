package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strings"
)

func response(conn net.Conn) {
	defer conn.Close()
	response := "HTTP/1.1 200 OK\r\n" +
		"Content-Type: text/plain\r\n" +
		"Content-Length: 13\r\n" +
		//"Connection: close\r\n" +
		"\r\n" +
		"Hello, World!"
	buffer := []byte(response)
	conn.Write(buffer)
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	//POST / HTTP/1.1
	schema, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println(err)
	}
	proc := strings.Split(schema, " ")
	method := proc[0]
	url := proc[1]
	httpVersion := proc[2]
	headers := make(map[string]string)

	fmt.Printf("接收到请求：%s %s %s\n", method, url, httpVersion)
	for {
		line, err := reader.ReadString('\n')
		//fmt.Println(line)
		if err != nil {
			if err != io.EOF {
				fmt.Println("Error reading from connection:", err)
			}
			break
		}

		if strings.TrimSpace(line) == "" {
			response(conn)
			break
		}

		header := strings.SplitN(line, ":", 2)
		//fmt.Println(header)
		headers[header[0]] = strings.TrimSpace(header[1])

	}

	for k, v := range headers {
		fmt.Printf("%s: %s\n", k, v)
	}

	//request := Request{URL: url, Headers: headers, Method: method}
}

func main() {
	listen, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}

	defer listen.Close()

	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		fmt.Println("Accepting connection from", conn.RemoteAddr())
		go handleConnection(conn) // Handle each connection in a separate goroutine
	}
}
