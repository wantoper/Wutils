package WTls

import (
	"WUtils/WTls/Const"
	"WUtils/WTls/Msg"
	"fmt"
	"net"
)

func Dial(address string) (*TlsConn, error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}

	tlsConn := newTlsConn(conn)
	tlsConn.handShakeFn = handShakeFunc
	return &tlsConn, nil
}

func handShakeFunc(conn net.Conn) (net.Conn, error) {
	cipherSuites := []uint8{Const.AES_GCM_128, Const.AES_GCM_256, Const.AES_GCM_128}
	hello := Msg.ClientHello{
		Version:      Const.Version1_1,
		SuiteLength:  uint8(len(cipherSuites)),
		CipherSuites: cipherSuites,
	}

	marshal := hello.Marshal()
	_, err := conn.Write(marshal)
	if err != nil {
		fmt.Print("Write error: ", err)
		return conn, err
	}

	serverHello := Msg.ServerHello{}
	bytes := make([]byte, 1024)
	n, err := conn.Read(bytes)
	if err != nil {
		fmt.Print("Read error: ", err)
	}
	serverHello.Unmarshal(bytes[:n])
	
	return conn, nil
}
