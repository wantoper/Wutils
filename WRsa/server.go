package WRsa

import (
	"WUtils/Util"
	"fmt"
	"net"
)

type RsaServer struct {
	listener net.Listener
}

func NewRsaServer(address string) (RsaServer, error) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return RsaServer{}, err
	}

	server := RsaServer{listener: listener}
	fmt.Printf("RsaServer started on %s\n", address)
	return server, nil
}

func (s RsaServer) Accept() (net.Conn, error) {
	conn, err := s.listener.Accept()
	if err != nil {
		return nil, err
	}

	publickey, _ := Util.GetPublicKey("server.crt")
	privateKey, _ := Util.GetPrivateKey("server.key")
	return newRsaConn(conn, publickey, privateKey), nil
}
