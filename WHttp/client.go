package WHttp

import (
	"WUtils/WTls"
	"bufio"
	"fmt"
	"net"
	"net/http"
	"net/url"
)

type HttpClient struct {
	Request *Request
	Conn    net.Conn
}

func NewHttpClient(reqUrl string, reqHeader Header) *HttpClient {
	uri, _ := url.ParseRequestURI(reqUrl)

	request := &Request{
		Url:    uri,
		Method: "GET",
		Proto:  "HTTP/1.1",
		Header: reqHeader,
	}

	return &HttpClient{
		Request: request,
	}
}

func (c *HttpClient) Do() {
	conn, _ := WTls.Dial(c.Request.Url.Host)
	writer := bufio.NewWriter(conn)
	fmt.Fprintf(writer, "%s %s %s\r\n", c.Request.Method, c.Request.Url.Path, c.Request.Proto)
	c.Request.Header.WriteHeaders(writer)
	writer.Flush()

	reader := bufio.NewReader(conn)

	//response := make([]byte, 4096)
	//n, err := conn.Read(response)
	//if err != nil {
	//	panic(err)
	//}
	//responseStr := string(response[:n])
	//println("Response from server:")
	//println(responseStr)
}

func HiHttp() {
	//uri, _ := url.ParseRequestURI("http://localhost:8080/hello")
	//fmt.Println(uri)
	//fmt.Println(uri.Host)
	//fmt.Println(uri.Scheme)
	//fmt.Println(uri.Path)
	//fmt.Println(uri.Port())
	//fmt.Println(uri.Query())
	//header := NewHeader()
	//header.Set("Content-Type", "text/html; charset=utf-8")
	//header.Set("ApiKey", "api-key-12345")

	//client := http.DefaultClient
	//request, err := http.NewRequest("GET", "http://localhost:8080/hello", nil)
	//client.Get()
	//client.Do(request)
	resp, err := http.Get("https://example.com")
	resp.Body.Close()
}
