package WHttp

import (
	"WUtils/WTls"
	"bufio"
	"fmt"
	"io"
	"net"
	"net/url"
	"strconv"
	"strings"
)

type HttpClient struct {
	Request *Request
	Conn    net.Conn
}

func (c *HttpClient) Do() *Response {
	conn, _ := WTls.Dial(c.Request.Url.Host)
	writer := bufio.NewWriter(conn)
	fmt.Fprintf(writer, "%s %s %s\r\n", c.Request.Method, c.Request.Url.Path, c.Request.Proto)

	var bodyData []byte
	if c.Request.Body != nil {
		bodyData, _ = io.ReadAll(c.Request.Body)
	}

	c.Request.Header.Set("Content-Length", strconv.Itoa(len(bodyData)))
	c.Request.Header.WriteHeaders(writer)

	if len(bodyData) > 0 {
		writer.Write(bodyData)
	}

	writer.Flush()
	return parserResponse(conn)
}

func Get(reqUrl string, reqHeader Header) *Response {
	uri, _ := url.ParseRequestURI(reqUrl)

	request := &Request{
		Url:    uri,
		Method: "GET",
		Proto:  "HTTP/1.1",
		Header: reqHeader,
	}

	client := HttpClient{
		Request: request,
	}

	return client.Do()
}

func Post(reqUrl string, reqHeader Header, body io.Reader) *Response {
	uri, _ := url.ParseRequestURI(reqUrl)

	request := &Request{
		Url:    uri,
		Method: "POST",
		Proto:  "HTTP/1.1",
		Header: reqHeader,
		Body:   body,
	}

	client := HttpClient{
		Request: request,
	}

	return client.Do()
}

func parserResponse(reader io.Reader) *Response {
	bfr := bufio.NewReader(reader)

	schema, _ := bfr.ReadString('\n')
	proto, rest, _ := strings.Cut(schema, " ")
	statusCode, _, _ := strings.Cut(rest, " ")
	statusCodeInt, _ := strconv.Atoi(statusCode)

	header := ParserHeader(bfr)

	contentLengthStr := header.Get("Content-Length")
	contentLength, err := strconv.ParseInt(contentLengthStr, 10, 64)
	if err != nil {
		fmt.Println("Error parsing Content-Length:", err)
		return nil
	}

	response := &Response{
		StatusCode:    statusCodeInt,
		Header:        header,
		Proto:         proto,
		ContentLength: contentLength,
		Body:          io.LimitReader(bfr, contentLength),
	}

	return response
}
