package WRsa

import "net"

type RsaClient struct {
	conn net.Conn
}

func NewRsaClient(address string) (*RsaClient, error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}

	publickey, _ := getPublicKey("server.crt")
	privateKey, _ := getPrivateKey("server.key")

	rsaConn := newRsaConn(conn, publickey, privateKey)
	client := &RsaClient{conn: rsaConn}
	return client, nil
}

func (c *RsaClient) Write(data []byte) (int, error) {
	n, err := c.conn.Write(data)
	if err != nil {
		return 0, err
	}
	return n, nil
}

func (c *RsaClient) Read(data []byte) (int, error) {
	n, err := c.conn.Read(data)
	if err != nil {
		return 0, err
	}
	return n, nil
}

func (c *RsaClient) Close() error {
	err := c.conn.Close()
	if err != nil {
		return err
	}
	return nil
}
