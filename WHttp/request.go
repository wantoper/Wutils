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
