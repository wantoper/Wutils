package WTls

import (
	"WUtils/WTls/Util"
	"fmt"
	"net"
	"time"
)

type HandShakeFunc func(tlsconn *TlsConn) error

type TlsConn struct {
	conn        net.Conn
	handshake   bool
	handShakeFn HandShakeFunc
	key         []byte
}

func newTlsConn(conn net.Conn) TlsConn {
	//key := []byte("12345678901234567891234567891234")
	return TlsConn{
		conn: conn,
		//key:  key,
	}
}

func (c *TlsConn) Read(b []byte) (n int, err error) {
	if c.HandShake() != nil {
		return 0, fmt.Errorf("handshake failed")
	}
	read, err := c.conn.Read(b)
	if err != nil {
		fmt.Printf("Read error: %v\n", err)
		return read, err
	}

	bytes, err := Util.Decrypt_AES(b[:read], c.key)
	fmt.Println(c.key)
	if err != nil {
		fmt.Printf("Decrypt error: %v\n", err)
	}
	copy(b, bytes)
	return len(bytes), err
}

func (c *TlsConn) Write(b []byte) (n int, err error) {
	if c.HandShake() != nil {
		return 0, fmt.Errorf("handshake failed")
	}
	encrypt, err := Util.Encrypt_AES(b, c.key)
	//encrypt, err := Util.Encrypt_Rsa(b, c.publicKey)
	if err != nil {
		fmt.Printf("Encrypt error: %v\n", err)
		return 0, err
	}
	return c.conn.Write(encrypt)
}

func (c *TlsConn) Close() error {
	return c.conn.Close()
}

func (c *TlsConn) LocalAddr() net.Addr {
	return c.conn.LocalAddr()
}

func (c *TlsConn) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

func (c *TlsConn) SetDeadline(t time.Time) error {
	return c.conn.SetDeadline(t)
}

func (c *TlsConn) SetReadDeadline(t time.Time) error {
	return c.conn.SetReadDeadline(t)
}

func (c *TlsConn) SetWriteDeadline(t time.Time) error {
	return c.conn.SetWriteDeadline(t)
}

func (c *TlsConn) HandShake() error {
	if !c.handshake {
		c.handshake = true
		fmt.Println("*[WTLS]WTLS handshake ....")
		err := c.handShakeFn(c)
		fmt.Printf("*[WTLS]WTLS handshake done,Session Key:%v\n", c.key)
		if err != nil {
			fmt.Printf("HandShake error: %v\n", err)
			return err
		}
	}
	return nil
}
