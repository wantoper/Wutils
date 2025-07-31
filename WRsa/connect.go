package WRsa

import (
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

	bytes, err := decrypt(b[:read], c.privateKey)
	copy(b, bytes)
	return len(bytes), err
}

func (c RsaConn) Write(b []byte) (n int, err error) {
	encrypt, err := encrypt(b, c.publicKey)
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
