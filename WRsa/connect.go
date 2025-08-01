package WRsa

import (
	"WUtils/Util"
	"crypto/rsa"
	"fmt"
	"net"
	"time"
)

type RsaConn struct {
	conn       net.Conn
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
}

func newRsaConn(conn net.Conn, publicKey *rsa.PublicKey, privateKey *rsa.PrivateKey) RsaConn {
	return RsaConn{
		conn:       conn,
		publicKey:  publicKey,
		privateKey: privateKey,
	}
}

func (c RsaConn) Read(b []byte) (n int, err error) {
	read, err := c.conn.Read(b)
	if err != nil {
		fmt.Printf("Read error: %v\n", err)
	}

	bytes, err := Util.Decrypt_Rsa(b[:read], c.privateKey)
	if err != nil {
		fmt.Printf("Decrypt error: %v\n", err)
	}
	copy(b, bytes)
	return len(bytes), err
}

func (c RsaConn) Write(b []byte) (n int, err error) {
	encrypt, err := Util.Encrypt_Rsa(b, c.publicKey)
	if err != nil {
		fmt.Printf("Encrypt error: %v\n", err)
		return 0, err
	}
	return c.conn.Write(encrypt)
}

func (c RsaConn) Close() error {
	return c.conn.Close()
}

func (c RsaConn) LocalAddr() net.Addr {
	return c.conn.LocalAddr()
}

func (c RsaConn) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

func (c RsaConn) SetDeadline(t time.Time) error {
	return c.conn.SetDeadline(t)
}

func (c RsaConn) SetReadDeadline(t time.Time) error {
	return c.conn.SetReadDeadline(t)
}

func (c RsaConn) SetWriteDeadline(t time.Time) error {
	return c.conn.SetWriteDeadline(t)
}
