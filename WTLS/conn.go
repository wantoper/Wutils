package WTLS

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"time"
)

type WTLSConn struct {
	conn     net.Conn
	cert     *x509.Certificate
	privKey  *rsa.PrivateKey
	aesKey   []byte
	aesIV    []byte
	cipher   cipher.AEAD
	isServer bool
}

func (c *WTLSConn) LocalAddr() net.Addr {
	//TODO implement me
	panic("implement me")
}

func (c *WTLSConn) RemoteAddr() net.Addr {
	//TODO implement me
	panic("implement me")
}

func (c *WTLSConn) SetDeadline(t time.Time) error {
	//TODO implement me
	panic("implement me")
}

func (c *WTLSConn) SetReadDeadline(t time.Time) error {
	//TODO implement me
	panic("implement me")
}

func (c *WTLSConn) SetWriteDeadline(t time.Time) error {
	//TODO implement me
	panic("implement me")
}

// NewWTLSConn 创建一个新的 WTLS 连接
func NewWTLSConn(conn net.Conn, cert *x509.Certificate, privKey *rsa.PrivateKey) *WTLSConn {
	return &WTLSConn{
		conn:     conn,
		cert:     cert,
		privKey:  privKey,
		isServer: privKey != nil,
	}
}

// Handshake 执行 TLS 握手
func (c *WTLSConn) Handshake() error {
	if c.isServer {
		return c.serverHandshake()
	}
	return c.clientHandshake()
}

func (c *WTLSConn) serverHandshake() error {
	// 1. 发送证书
	certBytes := c.cert.Raw
	if err := c.sendWithLength(certBytes); err != nil {
		return fmt.Errorf("发送证书失败: %v", err)
	}

	// 2. 接收客户端的 AES 密钥
	encryptedKey, err := c.receiveWithLength()
	if err != nil {
		return fmt.Errorf("接收密钥失败: %v", err)
	}

	// 3. 解密 AES 密钥
	c.aesKey, err = rsa.DecryptOAEP(sha256.New(), rand.Reader, c.privKey, encryptedKey, nil)
	if err != nil {
		return fmt.Errorf("解密密钥失败: %v", err)
	}

	// 4. 接收 IV
	c.aesIV, err = c.receiveWithLength()
	if err != nil {
		return fmt.Errorf("接收 IV 失败: %v", err)
	}

	// 5. 初始化 AES-GCM
	block, err := aes.NewCipher(c.aesKey)
	if err != nil {
		return fmt.Errorf("创建 AES 加密器失败: %v", err)
	}

	c.cipher, err = cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("创建 GCM 模式失败: %v", err)
	}

	return nil
}

func (c *WTLSConn) clientHandshake() error {
	// 1. 接收服务器证书
	certBytes, err := c.receiveWithLength()
	if err != nil {
		return fmt.Errorf("接收证书失败: %v", err)
	}

	cert, err := x509.ParseCertificate(certBytes)
	if err != nil {
		return fmt.Errorf("解析证书失败: %v", err)
	}
	c.cert = cert

	// 2. 生成并发送 AES 密钥
	c.aesKey = make([]byte, 32) // 256 位密钥
	if _, err := rand.Read(c.aesKey); err != nil {
		return fmt.Errorf("生成 AES 密钥失败: %v", err)
	}
	puk := c.cert.PublicKey.(rsa.PublicKey)
	encryptedKey, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, &puk, c.aesKey, nil)
	if err != nil {
		return fmt.Errorf("加密 AES 密钥失败: %v", err)
	}

	if err := c.sendWithLength(encryptedKey); err != nil {
		return fmt.Errorf("发送加密密钥失败: %v", err)
	}

	// 3. 生成并发送 IV
	c.aesIV = make([]byte, 12) // GCM 模式的 IV
	if _, err := rand.Read(c.aesIV); err != nil {
		return fmt.Errorf("生成 IV 失败: %v", err)
	}

	if err := c.sendWithLength(c.aesIV); err != nil {
		return fmt.Errorf("发送 IV 失败: %v", err)
	}

	// 4. 初始化 AES-GCM
	block, err := aes.NewCipher(c.aesKey)
	if err != nil {
		return fmt.Errorf("创建 AES 加密器失败: %v", err)
	}

	c.cipher, err = cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("创建 GCM 模式失败: %v", err)
	}

	return nil
}

// Read 实现 io.Reader 接口
func (c *WTLSConn) Read(b []byte) (n int, err error) {
	// 读取加密数据长度
	fmt.Println("来了 开始读取数据")

	buffer := make([]byte, 1024)

	//if _, err := io.ReadFull(c.conn, buffer); err != nil {
	//	return 0, err
	//}

	c.conn.Read(buffer)

	// 解密数据
	plaintext, err := c.cipher.Open(nil, c.aesIV, buffer, nil)
	if err != nil {
		return 0, err
	}

	return copy(b, plaintext), nil
}

// Write 实现 io.Writer 接口
func (c *WTLSConn) Write(b []byte) (n int, err error) {
	// 加密数据
	ciphertext := c.cipher.Seal(nil, c.aesIV, b, nil)

	// 写入加密数据长度
	if err := binary.Write(c.conn, binary.BigEndian, uint32(len(ciphertext))); err != nil {
		return 0, err
	}

	// 写入加密数据
	if _, err := c.conn.Write(ciphertext); err != nil {
		return 0, err
	}

	return len(b), nil
}

// Close 实现 io.Closer 接口
func (c *WTLSConn) Close() error {
	return c.conn.Close()
}

func (c *WTLSConn) sendWithLength(data []byte) error {
	if err := binary.Write(c.conn, binary.BigEndian, uint32(len(data))); err != nil {
		return err
	}
	_, err := c.conn.Write(data)
	return err
}

func (c *WTLSConn) receiveWithLength() ([]byte, error) {
	var length uint32
	if err := binary.Read(c.conn, binary.BigEndian, &length); err != nil {
		return nil, err
	}

	data := make([]byte, length)
	if _, err := io.ReadFull(c.conn, data); err != nil {
		return nil, err
	}

	return data, nil
}
