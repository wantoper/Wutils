package WWTls

import (
	"WUtils/WWTls/msg"
	"fmt"
	"net"
)

func Read(conn net.Conn, n int) []byte {
	bs := make([]byte, n)
	read, err := conn.Read(bs)
	if err != nil {
		panic(err)
	}
	if read != n {
		panic("read not enough")
	}
	return bs
}

func main() {
	listen, err := net.Listen("tcp", ":8443")
	if err != nil {
		panic(err)
	}
	defer listen.Close()

	for {
		conn, err := listen.Accept()
		if err != nil {
			panic(err)
		}

		recordHead := Read(conn, 5)

		contentType := msg.ContentType(recordHead[0])
		version := msg.TlsVersion(int(recordHead[1])<<8 | int(recordHead[2]))
		length := int(recordHead[3])<<8 | int(recordHead[4])
		fmt.Printf("Content Type: %v, Version: %x, Length: %v\n", contentType, version, length)

		tlsPainText := Read(conn, length)

		tlsPainTextHead := tlsPainText[0:4]
		handshakeType := msg.HandshakeType(tlsPainTextHead[0])
		handshakeLength := int(tlsPainTextHead[1])<<16 | int(tlsPainTextHead[2])<<8 | int(tlsPainTextHead[3])
		conn.Close()
	}
}
