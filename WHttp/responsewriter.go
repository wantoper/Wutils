package WHttp

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
)

type ResponseWriter interface {
	Header() Header
	Write([]byte) (int, error)
	SetStatus(statusCode int)
}

type ResponseWriterImpl struct {
	conn      *net.Conn
	HeaderMap Header
	w         bufio.Writer
}

func (rw *ResponseWriterImpl) Header() Header {
	if rw.HeaderMap == nil {
		rw.HeaderMap = make(Header)
	}
	return rw.HeaderMap
}

func (rw *ResponseWriterImpl) Write(b []byte) (int, error) {
	n, err := rw.w.Write(b)
	return n, err
}

func (rw *ResponseWriterImpl) SetStatus(statusCode int) {
	fmt.Fprintf(rw, "HTTP/1.1 %d %s\r\n", statusCode, http.StatusText(statusCode))
}
