package main

import "WUtils/WRsa"

func main() {
	client, err := WRsa.NewRsaClient(":4443")
	if err != nil {
		panic(err)
	}
	client.Write([]byte("hello world test rsa world test rsaworld test rsaworld test rsaworld test rsaworld test rsaworld test rsaworld test rsaworld test rsaworld test rsaworld test rsaworld test rsaworld test rsaworld test rsa"))
}
