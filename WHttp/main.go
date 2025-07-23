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

	// Read request line by line
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				fmt.Println("Error reading from connection:", err)
			}
			break
		}
		fmt.Println("Line:", strings.TrimSpace(line))

		// Check for the end of headers
		if strings.TrimSpace(line) == "" {
			response(conn)
			break
		}
	}
	fmt.Println("Connection closed")
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
