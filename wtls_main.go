package main

import (
	"WUtils/WTls"
	"WUtils/WTls/Util"
	"fmt"
)

func main() {
	publickey, _ := Util.GetPublicKey("server.crt")
	privateKey, _ := Util.GetPrivateKey("server.key")
	server, err := WTls.NewTlsServer(":4443", publickey, privateKey)
	if err != nil {
		panic(err)
	}

	for {
		conn, err := server.Accept()
		if err != nil {
			panic(err)
		}

		datas := make([]byte, 1024)
		n, err := conn.Read(datas)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(datas[:n]))
	}

	//WRsa.TestFun()
}
