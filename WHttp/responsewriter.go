package WHttp

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"strconv"
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

func (rw *ResponseWriterImpl) SetStatus(statusCode int) {
	if rw.writeStatus {
		return
	}
	fmt.Fprintf(&rw.w, "HTTP/1.1 %d %s\r\n", statusCode, http.StatusText(statusCode))
	rw.writeStatus = true
}

func (rw *ResponseWriterImpl) writeHeader() {
	if !rw.writeStatus {
		rw.SetStatus(200)
	}
	if rw.HeaderMap == nil {
		return
	}
	rw.HeaderMap.WriteHeaders(&rw.w)
}

func (rw *ResponseWriterImpl) Write(b []byte) (int, error) {
	rw.Header().Set("Content-Length", strconv.Itoa(len(b)))
	rw.writeHeader()
	n, err := rw.w.Write(b)
	return n, err
}

func (rw *ResponseWriterImpl) FinishRequest() {
	rw.w.Flush()
}
