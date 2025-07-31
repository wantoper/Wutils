package WTLS

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net"
	"os"
)

type WTLSClient struct {
	conn *WTLSConn
	cert *x509.Certificate
}

// NewWTLSClient 创建一个新的 WTLS 客户端
func NewWTLSClient() (*WTLSClient, error) {
	client := &WTLSClient{}

	// 尝试读取服务器证书
	if _, err := os.Stat("server.crt"); err == nil {
		certPEM, err := os.ReadFile("server.crt")
		if err != nil {
			return nil, fmt.Errorf("读取证书文件失败: %v", err)
		}

		// 解析证书
		block, _ := pem.Decode(certPEM)
		if block == nil {
			return nil, fmt.Errorf("解析证书PEM失败")
		}

		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("解析证书失败: %v", err)
		}

		client.cert = cert
	}

	return client, nil
}

// Connect 连接到指定地址的服务器
func (c *WTLSClient) Connect(address string) error {
	// 建立 TCP 连接
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return fmt.Errorf("连接服务器失败: %v", err)
	}

	// 创建 WTLS 连接，如果有预加载的证书就使用它
	c.conn = NewWTLSConn(conn, c.cert, nil)
	fmt.Println("连接到服务器:", address)
	// 执行握手
	if err := c.conn.Handshake(); err != nil {
		conn.Close()
		return fmt.Errorf("TLS 握手失败: %v", err)
	}

	return nil
}

// Write 写入数据
func (c *WTLSClient) Write(data []byte) (int, error) {
	return c.conn.Write(data)
}

// Read 读取数据
func (c *WTLSClient) Read(data []byte) (int, error) {
	return c.conn.Read(data)
}

// Close 关闭连接
func (c *WTLSClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
