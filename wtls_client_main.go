package main

import "WUtils/WRsa"

func main() {
	client, err := WRsa.NewRsaClient(":4443")
	if err != nil {
		panic(err)
	}
	client.Write([]byte("测试的问题啊啊啊测试的问题啊啊啊测试的问题啊啊啊测试的问题啊啊啊测试的问题啊啊啊测试的问题啊题啊啊啊测试的问题啊题啊啊啊测试的问题啊"))
}
