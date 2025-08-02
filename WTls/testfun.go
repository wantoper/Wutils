package WTls

import (
	"WUtils/WTls/Const"
	"WUtils/WTls/Msg"
	"fmt"
)

func TestFun() string {
	cipherSuites := []uint8{Const.AES_GCM_128, Const.AES_GCM_256, Const.AES_GCM_128, Const.AES_GCM_128, Const.AES_GCM_128}
	hello := Msg.ClientHello{Version: Const.Version1_1, SuiteLength: uint8(255), CipherSuites: cipherSuites}

	marshal := hello.Marshal()
	fmt.Println(marshal)
	fmt.Printf("%x\n", marshal)
	fmt.Println(string(marshal))

	hello2 := Msg.ClientHello{}
	hello2.UnmarShal(marshal)
	fmt.Println(hello2.Version)
	fmt.Println(hello2.SuiteLength)
	fmt.Println(hello2.CipherSuites)

	//serverHello := Msg.ServerHello{
	//	Version:     Msg.Version1_1,
	//	CipherSuite: Msg.AES_GCM_128,
	//	KeyLength:   16,
	//	EncryptKey:  []byte("1234567890123456"),
	//}
	//
	//marshal := serverHello.Marshal()
	//fmt.Printf("%x\n", marshal)
	//
	//serverHello2 := Msg.ServerHello{}
	//
	//err := serverHello2.Unmarshal(marshal)
	//if err != nil {
	//	fmt.Println("Unmarshal error:", err)
	//	return ""
	//}
	//
	//fmt.Println("Version:", serverHello2.Version)
	//fmt.Println("CipherSuite:", serverHello2.CipherSuite)
	//fmt.Println("KeyLength:", serverHello2.KeyLength)
	//fmt.Println("EncryptKey:", string(serverHello2.EncryptKey))
	//
	return ""
}
