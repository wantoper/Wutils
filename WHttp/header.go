package WHttp

import (
	"fmt"
	"io"
)

type Header map[string]string

func NewHeader() Header {
	return make(Header)
}

func (h Header) Set(key, value string) {
	h[key] = value
}

func (h Header) Get(key string) string {
	return h[key]
}

func (h Header) WriteHeaders(w io.Writer) {
	for k, v := range h {
		fmt.Fprintf(w, "%s: %s\r\n", k, v)
	}
	fmt.Fprint(w, "\r\n")
}
