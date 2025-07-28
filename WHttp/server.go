package WHttp

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/url"
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

func create_request(conn net.Conn) (*Request, error) {
	reader := bufio.NewReader(conn)

	schema, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println(err)
	}
	mehtod, rest, _ := strings.Cut(schema, " ")
	requestURI, proto, _ := strings.Cut(rest, " ")
	rawURI := "http://" + requestURI
	parseRequestURI, _ := url.ParseRequestURI(rawURI)

	headers := make(Header)
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

		headerK, headerV, _ := strings.Cut(line, ":")
		headers[headerK] = strings.TrimSpace(headerV)
	}

	if parseRequestURI.Host == "" {
		parseRequestURI.Host = headers.Get("Host")
	}
	contenlength := headers.Get("Content-Length")
	realLength, err := strconv.ParseInt(contenlength, 10, 64)
	req := &Request{
		Method: mehtod,
		Proto:  proto,
		Url:    parseRequestURI,
		Header: headers,
		//Body:   reader,
		Body: io.LimitReader(reader, realLength),
	}

	return req, err
}

func handleConnection(conn net.Conn) {

	request, _ := create_request(conn)
	fmt.Println("New connection from", request.Url)

	//for k, v := range request.Header {
	//	fmt.Printf("%s: %s\n", k, v)
	//}

	//contentlen := request.Header.Get("Content-Length")
	//contentLength, _ := strconv.Atoi(contentlen)
	//
	//body := make([]byte, contentLength)
	//request.Body.Read(body)
	//
	//bodystr := string(body)
	//fmt.Println(bodystr)

	all, _ := io.ReadAll(request.Body)
	fmt.Println("==================")
	fmt.Println("stringall:", len(all))
	bodystr := string(all)
	fmt.Println(bodystr)

	response(conn)

	conn.Close()
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
