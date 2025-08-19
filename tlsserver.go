package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"io"
	"math/big"
	"net"
	"os"
	"time"
)

// LoggerConn is a wrapper around net.Conn that logs all read and write operations in hex dump format.
// This helps visualize the TLS handshake messages as they are sent and received over the wire.
type LoggerConn struct {
	net.Conn
	prefix string // e.g., "Server" or "Client" to distinguish logs
}

func (c *LoggerConn) Read(b []byte) (int, error) {
	n, err := c.Conn.Read(b)
	if n > 0 {
		fmt.Printf("%s Received:\n%s\n", c.prefix, hex.Dump(b[:n]))
	}
	return n, err
}

func (c *LoggerConn) Write(b []byte) (int, error) {
	if len(b) > 0 {
		fmt.Printf("%s Sent:\n%s\n", c.prefix, hex.Dump(b))
	}
	n, err := c.Conn.Write(b)
	return n, err
}

// generateSelfSignedCert generates a self-signed certificate for testing purposes.
// Returns the cert PEM and key PEM as byte slices.
func generateSelfSignedCert() ([]byte, []byte, error) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}

	notBefore := time.Now()
	notAfter := notBefore.Add(365 * 24 * time.Hour)

	serialNumber, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return nil, nil, err
	}

	template := x509.Certificate{
		SerialNumber:          serialNumber,
		Subject:               pkix.Name{Organization: []string{"Test Org"}},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		DNSNames:              []string{"localhost"},
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, priv.Public(), priv)
	if err != nil {
		return nil, nil, err
	}

	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})

	return certPEM, keyPEM, nil
}

func runServer() {
	// Generate self-signed cert for the server
	certPEM, keyPEM, err := generateSelfSignedCert()
	if err != nil {
		fmt.Printf("Error generating cert: %v\n", err)
		return
	}

	// Load the cert
	cert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		fmt.Printf("Error loading cert: %v\n", err)
		return
	}

	// TLS config forcing TLS 1.3
	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS13,
		MaxVersion:   tls.VersionTLS13,
	}

	// Listen on TCP
	listener, err := net.Listen("tcp", "localhost:8443")
	if err != nil {
		fmt.Printf("Error listening: %v\n", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server listening on localhost:8443")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Error accepting: %v\n", err)
			continue
		}

		// Wrap with LoggerConn to log handshake data
		loggerConn := &LoggerConn{Conn: conn, prefix: "Server"}

		tlsConn := tls.Server(loggerConn, config)

		// Perform handshake (this is where TLS 1.3 handshake happens, and logs will show sent/received data)
		if err := tlsConn.Handshake(); err != nil {
			fmt.Printf("Handshake error: %v\n", err)
			continue
		}

		// After handshake, handle the connection (e.g., read a message)
		buf := make([]byte, 1024)
		n, err := tlsConn.Read(buf)
		if err != nil && err != io.EOF {
			fmt.Printf("Read error: %v\n", err)
		} else {
			fmt.Printf("Server received message: %s\n", string(buf[:n]))
		}

		// Send a response
		_, err = tlsConn.Write([]byte("Hello from TLS 1.3 Server!"))
		if err != nil {
			fmt.Printf("Write error: %v\n", err)
		}

		tlsConn.Close()
	}
}

func runClient() {
	// TLS config forcing TLS 1.3, and skipping verify for self-signed cert (for testing)
	config := &tls.Config{
		InsecureSkipVerify: true,
		MinVersion:         tls.VersionTLS13,
		MaxVersion:         tls.VersionTLS13,
	}

	// Dial TCP
	conn, err := net.Dial("tcp", "localhost:8443")
	if err != nil {
		fmt.Printf("Error dialing: %v\n", err)
		return
	}

	// Wrap with LoggerConn to log handshake data
	loggerConn := &LoggerConn{Conn: conn, prefix: "Client"}

	tlsConn := tls.Client(loggerConn, config)

	// Perform handshake (this is where TLS 1.3 handshake happens, and logs will show sent/received data)
	if err := tlsConn.Handshake(); err != nil {
		fmt.Printf("Handshake error: %v\n", err)
		return
	}

	// Send a message
	_, err = tlsConn.Write([]byte("Hello from TLS 1.3 Client!"))
	if err != nil {
		fmt.Printf("Write error: %v\n", err)
		return
	}

	// Read response
	buf := make([]byte, 1024)
	n, err := tlsConn.Read(buf)
	if err != nil && err != io.EOF {
		fmt.Printf("Read error: %v\n", err)
	} else {
		fmt.Printf("Client received message: %s\n", string(buf[:n]))
	}

	tlsConn.Close()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go [server|client]")
		return
	}

	mode := os.Args[1]
	if mode == "server" {
		runServer()
	} else if mode == "client" {
		runClient()
	} else {
		fmt.Println("Invalid mode. Use 'server' or 'client'")
	}
}
