package WHttp

import "io"

type Response struct {
	Status        string
	StatusCode    int
	Proto         string
	Header        Header
	ContentLength int64
	Body          io.ReadCloser
}
