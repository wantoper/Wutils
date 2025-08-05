package WHttp

import (
	"io"
	"net/url"
)

type Request struct {
	Method string // HTTP method (GET, POST, etc.)
	Url    *url.URL
	Proto  string
	Header Header
	Body   io.Reader
}

type Response struct {
	StatusCode    int
	Proto         string
	Header        Header
	ContentLength int64
	Body          io.Reader
}
