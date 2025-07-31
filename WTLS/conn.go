package WTLS

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/binary"
	"encoding/pem"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
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

	var dataLen uint32

	err = binary.Read(c.conn, binary.LittleEndian, &dataLen)
	fmt.Println("读取数据长度:", dataLen)
	data := make([]byte, dataLen)
	if _, err := io.ReadFull(c.conn, data); err != nil {
		return 0, err
	}

	// 解密数据
	//plaintext, err := c.cipher.Open(nil, c.aesIV, data, nil)
	//if err != nil {
	//	return 0, err
	//}

	plaintext, err := decrypt(data, c.privKey)

	return copy(b, plaintext), nil
}

// Write 实现 io.Writer 接口
func (c *WTLSConn) Write(b []byte) (n int, err error) {
	// 加密数据
	//ciphertext := c.cipher.Seal(nil, c.aesIV, b, nil)
	bytes, err := encrypt(b, c.cert.PublicKey.(*rsa.PublicKey))
	ciphertext := bytes

	// 写入加密数据长度
	//if err := binary.Write(c.conn, binary.BigEndian, uint32(len(ciphertext))); err != nil {
	//	return 0, err
	//}

	// 写入加密数据
	if _, err := c.conn.Write(ciphertext); err != nil {
		return 0, err
	}
	fmt.Printf("发送数据: %s\n", b)
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

// 从证书文件读取公钥
func getPublicKey(certPath string) (*rsa.PublicKey, error) {
	certBytes, err := os.ReadFile(certPath)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(certBytes)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, err
	}
	return cert.PublicKey.(*rsa.PublicKey), nil
}

// 从私钥文件读取私钥
func getPrivateKey(keyPath string) (*rsa.PrivateKey, error) {
	keyBytes, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(keyBytes)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}
	return x509.ParsePKCS1PrivateKey(block.Bytes)
}

// 加密函数
func encrypt(plaintext []byte, publicKey *rsa.PublicKey) ([]byte, error) {
	return rsa.EncryptPKCS1v15(rand.Reader, publicKey, plaintext)
}

// 解密函数
func decrypt(ciphertext []byte, privateKey *rsa.PrivateKey) ([]byte, error) {
	return rsa.DecryptPKCS1v15(rand.Reader, privateKey, ciphertext)
}
