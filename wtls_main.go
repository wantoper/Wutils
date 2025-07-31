package main

import (
	"WUtils/WRsa"
	"fmt"
)

func main() {
	server, err := WRsa.NewRsaServer(":4443")
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
}
