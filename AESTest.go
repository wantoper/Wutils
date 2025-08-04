package main

import (
	"WUtils/WTls/Util"
	"fmt"
)

func main() {
	//randonKey, _ := getRandonKey()
	////randonKey := []byte("12345678901234567890123456789012")
	//fmt.Println("randonKey:", randonKey)
	//key, _ := Util.GetPublicKey("server.crt")
	//bytes, _ := Util.Encrypt_Rsa(randonKey, key)
	//
	//fmt.Println(bytes)
	//
	////encryptedData, err := encrypt_aes(rsaPBK, randonKey)
	////if err != nil {
	////	panic(err)
	////}
	////
	////fmt.Println("encryptedData:", string(encryptedData))
	////
	////decryptedData, _ := decrypt_aes(encryptedData, randonKey)
	////fmt.Println("decryptedData:", string(decryptedData))
	//
	//pri_Key, _ := Util.GetPrivateKey("server.key")
	//rsa, _ := Util.Decrypt_Rsa(bytes, pri_Key)
	//fmt.Println(rsa)

	key, _ := Util.GetPublicKey("server.crt")
	fmt.Println(key.E)
	fmt.Println(key.N)
	bytes, _ := Util.PublickeyToBytes(key)
	fmt.Println(len(bytes))
	fmt.Println("Public Key Bytes:", bytes)
	publicKey, _ := Util.BytesToPublicKey(bytes)
	fmt.Println(publicKey.E)
	fmt.Println(publicKey.N)
}
