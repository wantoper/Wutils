package WTls

import (
	"WUtils/WTls/Const"
	"WUtils/WTls/Msg"
	"WUtils/WTls/Util"
	"crypto/rsa"
	"fmt"
	"net"
)

type TlsServer struct {
	listener   net.Listener
	publickey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
}

func NewTlsServer(address string, publickey *rsa.PublicKey, privateKey *rsa.PrivateKey) (TlsServer, error) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return TlsServer{}, err
	}

	server := TlsServer{listener: listener}
	fmt.Printf("RsaServer started on %s\n", address)
	server.publickey = publickey
	server.privateKey = privateKey

	return server, nil
}

func (s TlsServer) Accept() (net.Conn, error) {
	conn, err := s.listener.Accept()
	if err != nil {
		return nil, err
	}

	//publickey, _ := Util.GetPublicKey("server.crt")
	//privateKey, _ := Util.GetPrivateKey("server.key")
	tlsConn := newTlsConn(conn)
	tlsConn.handShakeFn = s.HandShakeFunc
	return tlsConn, nil
}

func (s TlsServer) HandShakeFunc(conn net.Conn) (net.Conn, error) {
	client_hello := Msg.ClientHello{}
	bytes := make([]byte, 1024)
	n, _ := conn.Read(bytes)
	client_hello.UnmarShal(bytes[:n])

	use_cipher := client_hello.CipherSuites[0]
	key, _ := Util.GetRandonKey(use_cipher)
	encryptRsa, _ := Util.Encrypt_RSA(key, s.publickey)

	server_hello := Msg.ServerHello{
		Version:     client_hello.Version,
		CipherSuite: use_cipher,
		KeyLength:   uint8(len(key)),
		EncryptKey:  encryptRsa,
	}
	fmt.Printf("交换AES密钥 version: %v, cipherSuite: %v, keyLength: %v, key: %v, encryptKey: %v\n", server_hello.Version, Const.GetCipherSuiteName(server_hello.CipherSuite), server_hello.KeyLength, key, server_hello.EncryptKey)
	server_hello_data := server_hello.Marshal()
	conn.Write(server_hello_data)
	return conn, nil
}
