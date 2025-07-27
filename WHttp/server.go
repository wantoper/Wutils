package WHttp

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
)

func response(conn net.Conn) {
	defer conn.Close()
	writer := ResponseWriterImpl{
		w: *bufio.NewWriter(conn),
	}

	writer.SetStatus(302)
	response := "Content-Type: text/plain\r\n" +
		"Content-Length: 13\r\n" +
		"\r\n" +
		"Hello, World!"
	writer.w.WriteString(response)
	writer.w.Flush()
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
		if err != nil {
			if err != io.EOF {
				fmt.Println("Error reading from connection:", err)
			}
			break
		}

		if strings.EqualFold(line, "\r\n") {
			break
		}

		header := strings.SplitN(line, ":", 2)
		headers[header[0]] = strings.TrimSpace(header[1])
	}
	//
	for k, v := range headers {
		fmt.Printf("%s: %s\n", k, v)
	}

	_, ok := headers["Content-Length"]
	if !ok {
		fmt.Println("No Content-Length header")
	}
	contentLength, err := strconv.Atoi(headers["Content-Length"])
	if err != nil {
		fmt.Println("Invalid Content-Length header")
	}

	fmt.Println("内容长度：", contentLength)

	body := make([]byte, contentLength)
	_, err = reader.Read(body)
	if err != nil {
		fmt.Println("Error reading from connection:", err)
	}

	bodystr := string(body)
	fmt.Println(bodystr)

	//request := Request{URL: url, Headers: headers, Method: method}
	response(conn)
}

func StartServer(address string) error {
	listen, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Println("Error starting server:", err)
		return err
	}

	defer listen.Close()

	for {
		conn, err := listen.Accept()
		if err != nil {
			return err
		}
		fmt.Println("Accepting connection from", conn.RemoteAddr())
		go handleConnection(conn)
	}
}
