package WTls

import (
	"WUtils/WTls/Util"
	"fmt"
	"net"
	"time"
)

type HandShakeFunc func(conn net.Conn) (net.Conn, error)

type TlsConn struct {
	conn        net.Conn
	handshake   bool
	handShakeFn HandShakeFunc
	key         []byte
}

func newTlsConn(conn net.Conn) TlsConn {
	key := []byte("12345678901234567891234567891234")
	return TlsConn{
		conn: conn,
		key:  key,
	}
}

func (c TlsConn) Read(b []byte) (n int, err error) {
	if !c.handshake {
		c.handshake = true
		_, err = c.handShakeFn(c.conn)
		if err != nil {
			fmt.Printf("HandShake error: %v\n", err)
			return 0, err
		}
	}
	read, err := c.conn.Read(b)
	if err != nil {
		fmt.Printf("Read error: %v\n", err)
	}

	bytes, err := Util.Decrypt_AES(b[:read], c.key)
	//bytes, err := Util.Decrypt_Rsa(b[:read], c.privateKey)
	if err != nil {
		fmt.Printf("Decrypt error: %v\n", err)
	}
	copy(b, bytes)
	return len(bytes), err
}

func (c TlsConn) Write(b []byte) (n int, err error) {
	if !c.handshake {
		c.handshake = true
		fmt.Printf("pointer address: %p\n", c.handShakeFn)
		_, err = c.handShakeFn(c.conn)
		if err != nil {
			fmt.Printf("HandShake error: %v\n", err)
			return 0, err
		}
	}
	encrypt, err := Util.Encrypt_AES(b, c.key)
	//encrypt, err := Util.Encrypt_Rsa(b, c.publicKey)
	if err != nil {
		fmt.Printf("Encrypt error: %v\n", err)
		return 0, err
	}
	return c.conn.Write(encrypt)
}

func (c TlsConn) Close() error {
	return c.conn.Close()
}

func (c TlsConn) LocalAddr() net.Addr {
	return c.conn.LocalAddr()
}

func (c TlsConn) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

func (c TlsConn) SetDeadline(t time.Time) error {
	return c.conn.SetDeadline(t)
}

func (c TlsConn) SetReadDeadline(t time.Time) error {
	return c.conn.SetReadDeadline(t)
}

func (c TlsConn) SetWriteDeadline(t time.Time) error {
	return c.conn.SetWriteDeadline(t)
}
