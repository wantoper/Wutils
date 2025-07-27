package WHttp

import "io"

type Request struct {
	Method string // HTTP method (GET, POST, etc.)
	Url    string // Request URL
	Header Header
	Body   io.ReadCloser
}
