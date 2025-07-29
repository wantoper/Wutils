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
	conn        *net.Conn
	HeaderMap   Header
	w           bufio.Writer
	writeStatus bool
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
	if rw.writeStatus {
		return
	}
	fmt.Fprintf(rw, "HTTP/1.1 %d %s\r\n", statusCode, http.StatusText(statusCode))
	rw.writeStatus = true
}

func (rw *ResponseWriterImpl) WriteHeader() {
	if rw.HeaderMap == nil {
		return
	}
	for k, v := range rw.HeaderMap {
		fmt.Fprintf(&rw.w, "%s: %s\r\n", k, v)
	}
}

func (rw *ResponseWriterImpl) FinishRequest() {
	if rw.conn != nil {
		rw.WriteHeader()
		rw.w.Flush()
		(*rw.conn).Close()
		rw.conn = nil
	}
}
