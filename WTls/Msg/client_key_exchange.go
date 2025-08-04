package Msg

import (
	"errors"
	"unsafe"
)

type ClientKeyExchange struct {
	CipherSuite uint8
	KeyLength   uint16
	EncryptKey  []byte
}

func (c *ClientKeyExchange) Marshal() []byte {
	dataSize := int(unsafe.Sizeof(c.CipherSuite)) + int(unsafe.Sizeof(c.KeyLength)) + len(c.EncryptKey)
	data := make([]byte, dataSize)
	data[0] = c.CipherSuite
	data[1] = byte(c.KeyLength >> 8)
	data[2] = byte(c.KeyLength & 0xFF)
	copy(data[3:], c.EncryptKey)
	return data
}

func (c *ClientKeyExchange) Unmarshal(data []byte) error {
	if len(data) < 3 {
		return errors.New("data too short")
	}

	c.CipherSuite = data[0]
	c.KeyLength = uint16(data[1])<<8 | uint16(data[2])
	if len(data) < int(3+c.KeyLength) {
		return errors.New("data length does not match key length")
	}
	c.EncryptKey = make([]byte, c.KeyLength)
	copy(c.EncryptKey, data[3:3+c.KeyLength])
	return nil
}
