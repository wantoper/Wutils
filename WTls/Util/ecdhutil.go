package Util

import (
	"crypto/ecdh"
	"crypto/rand"
)

func GetECDHKey() (*ecdh.PrivateKey, *ecdh.PublicKey, error) {
	privateKey, err := ecdh.P256().GenerateKey(rand.Reader)
	if err != nil {
		return nil, nil, err
	}

	publicKey := privateKey.PublicKey()

	return privateKey, publicKey, nil
}
