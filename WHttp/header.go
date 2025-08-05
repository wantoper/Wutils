package WHttp

import (
	"bufio"
	"fmt"
	"io"
	"strings"
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

func ParserHeader(reader io.Reader) Header {
	header := NewHeader()

	bfr := reader.(*bufio.Reader)

	for {
		line, err := bfr.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				fmt.Println("Error reading from connection:", err)
			}
			break
		}

		if strings.EqualFold(line, "\r\n") {
			break
		}

		headerK, headerV, _ := strings.Cut(line, ":")
		header[headerK] = strings.TrimSpace(headerV)
	}

	return header
}
