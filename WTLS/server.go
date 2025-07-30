package WTLS

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"os"
	"time"
)

type WTLSServer struct {
	cert     *x509.Certificate
	privKey  *rsa.PrivateKey
	listener net.Listener
}

// NewWTLSServer 创建一个新的 WTLS 服务器
func NewWTLSServer() (*WTLSServer, error) {
	// 生成私钥
	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, fmt.Errorf("生成私钥失败: %v", err)
	}

	// 创建证书模板
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName:   "WTLS Server",
			Organization: []string{"WTLS"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(time.Hour * 24 * 365), // 1年有效期
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	// 创建自签名证书
	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &privKey.PublicKey, privKey)
	if err != nil {
		return nil, fmt.Errorf("创建证书失败: %v", err)
	}

	cert, err := x509.ParseCertificate(derBytes)
	if err != nil {
		return nil, fmt.Errorf("解析证书失败: %v", err)
	}

	// 保存证书到文件
	certOut, err := os.Create("server.crt")
	if err != nil {
		return nil, fmt.Errorf("创建证书文件失败: %v", err)
	}
	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	certOut.Close()

	// 保存私钥到文件
	keyOut, err := os.Create("server.key")
	if err != nil {
		return nil, fmt.Errorf("创建私钥文件失败: %v", err)
	}
	pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privKey)})
	keyOut.Close()

	return &WTLSServer{
		cert:    cert,
		privKey: privKey,
	}, nil
}

func (s *WTLSServer) Listen(address string) error {
	var err error
	s.listener, err = net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("监听失败: %v", err)
	}
	return nil
}

func (s *WTLSServer) Accept() (net.Conn, error) {
	conn, err := s.listener.Accept()
	if err != nil {
		return nil, fmt.Errorf("接受连接失败: %v", err)
	}

	// 包装原始连接为 WTLS 连接
	return NewWTLSConn(conn, s.cert, s.privKey), nil
}

func (s *WTLSServer) Close() error {
	if s.listener != nil {
		return s.listener.Close()
	}
	return nil
}
