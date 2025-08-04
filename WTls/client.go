package WTls

import (
	"WUtils/WTls/Msg"
	"WUtils/WTls/Util"
	"WUtils/WTls/consts"
	"crypto/rsa"
	"fmt"
	"net"
)

func Dial(address string) (*TlsConn, error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}

	client := TlsClient{}

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
		Version:      consts.Version1_0,
		SuiteLength:  uint8(len(cipherSuites)),
		CipherSuites: cipherSuites,
	}

	marshal := hello.Marshal()
	_, err := tlsconn.conn.Write(marshal)
	if err != nil {
		fmt.Print("Write error: ", err)
		return err
	}

	server_hello := Msg.ServerHello{}
	bytes := make([]byte, 1024)
	n, err := tlsconn.conn.Read(bytes)
	if err != nil {
		fmt.Print("Read error: ", err)
	}
	server_hello.Unmarshal(bytes[:n])
	publickey, err := Util.BytesToPublicKey(server_hello.EncryptKey)
	fmt.Printf("*[WTLS]接收公钥 version: %v, cipherSuite: %v, keyLength: %v, key: %v, encryptKey: %v\n", server_hello.Version, consts.GetCipherSuiteName(server_hello.CipherSuite), server_hello.KeyLength, publickey.E, server_hello.EncryptKey)

	key, err := Util.GetRandonKey(server_hello.CipherSuite)
	tlsconn.key = key
	fmt.Println("*[WTLS]选择的密钥:", key)
	key, _ = Util.Encrypt_RSA(key, publickey)
	if err != nil {
		fmt.Print("Get random key error: ", err)
		return err
	}

	exchange := Msg.ClientKeyExchange{
		CipherSuite: server_hello.CipherSuite,
		KeyLength:   uint16(len(key)),
		EncryptKey:  key,
	}

	exchange.Marshal()
	_, err = tlsconn.conn.Write(exchange.Marshal())
	return nil
}
