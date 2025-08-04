package WTls

import (
	"WUtils/WTls/Msg"
	"WUtils/WTls/Util"
	"WUtils/WTls/consts"
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

	tlsConn := newTlsConn(conn)
	tlsConn.handShakeFn = s.HandShakeFunc
	return &tlsConn, nil
}

func (s TlsServer) HandShakeFunc(tlsconn *TlsConn) error {
	client_hello := Msg.ClientHello{}
	bytes := make([]byte, 1024)
	n, _ := tlsconn.conn.Read(bytes)
	client_hello.UnmarShal(bytes[:n])

	use_cipher := client_hello.CipherSuites[0]
	keyBytes, _ := Util.PublickeyToBytes(s.publickey)

	server_hello := Msg.ServerHello{
		Version:     client_hello.Version,
		CipherSuite: use_cipher,
		KeyLength:   uint16(len(keyBytes)),
		EncryptKey:  keyBytes,
	}
	fmt.Printf("*[WTLS]发送公钥 version: %v, cipherSuite: %v, keyLength: %v, key: %v, encryptKey: %v\n", server_hello.Version, consts.GetCipherSuiteName(server_hello.CipherSuite), server_hello.KeyLength, s.publickey.E, server_hello.EncryptKey)
	server_hello_data := server_hello.Marshal()
	tlsconn.conn.Write(server_hello_data)

	n, _ = tlsconn.conn.Read(bytes)
	exchange := Msg.ClientKeyExchange{}
	err := exchange.Unmarshal(bytes[:n])
	if err != nil {
		fmt.Println("Unmarshal error:", err)
		return err
	}
	decryptRsa, err := Util.Decrypt_RSA(exchange.EncryptKey, s.privateKey)

	fmt.Printf("*[WTLS]接收密钥 exchange:cipherSuite: %v, keyLength: %v, encryptKey: %v\n", consts.GetCipherSuiteName(exchange.CipherSuite), exchange.KeyLength, decryptRsa)
	tlsconn.key = decryptRsa
	return nil
}
