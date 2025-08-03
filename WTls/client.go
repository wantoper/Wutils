package WTls

import (
	"WUtils/WTls/Msg"
	"WUtils/WTls/consts"
	"crypto/rsa"
	"fmt"
	"net"
)

func Dial(address string, publickey *rsa.PublicKey) (*TlsConn, error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}

	client := TlsClient{
		publickey: publickey,
	}

	tlsConn := newTlsConn(conn)
	tlsConn.handShakeFn = client.HandShakeFunc
	return &tlsConn, nil
}

type TlsClient struct {
	publickey *rsa.PublicKey
}

func (tc *TlsClient) HandShakeFunc(tlsconn *TlsConn) error {
	cipherSuites := []uint8{consts.AES_GCM_128, consts.AES_GCM_256, consts.AES_GCM_128}
	hello := Msg.ClientHello{
		Version:      consts.Version1_1,
		SuiteLength:  uint8(len(cipherSuites)),
		CipherSuites: cipherSuites,
	}

	marshal := hello.Marshal()
	_, err := tlsconn.conn.Write(marshal)
	if err != nil {
		fmt.Print("Write error: ", err)
		return err
	}

	serverHello := Msg.ServerHello{}
	bytes := make([]byte, 1024)
	n, err := tlsconn.conn.Read(bytes)
	if err != nil {
		fmt.Print("Read error: ", err)
	}
	serverHello.Unmarshal(bytes[:n])
	//Util.Decrypt_RSA(serverHello.EncryptKey, tc.publickey)
	//待RSA解密
	tlsconn.key = serverHello.EncryptKey

	return nil
}
