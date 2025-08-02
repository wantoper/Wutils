package Msg

import (
	"errors"
	"unsafe"
)

type ServerHello struct {
	Version     uint16
	CipherSuite uint8
	KeyLength   uint8
	EncryptKey  []byte
}

func (s *ServerHello) Marshal() []byte {
	len := int(unsafe.Sizeof(s.Version)) + int(unsafe.Sizeof(s.CipherSuite)) + int(unsafe.Sizeof(s.KeyLength)) + len(s.EncryptKey)
	data := make([]byte, len)

	data[0] = byte(s.Version >> 8)
	data[1] = byte(s.Version & 0xFF)
	data[2] = s.CipherSuite
	data[3] = s.KeyLength
	copy(data[4:], s.EncryptKey)
	return data
}

func (s *ServerHello) Unmarshal(data []byte) error {
	if len(data) < 4 {
		return errors.New("data too short")
	}

	s.Version = uint16(data[0])<<8 | uint16(data[1])
	s.CipherSuite = data[2]
	s.KeyLength = data[3]
	s.EncryptKey = make([]byte, s.KeyLength)
	copy(s.EncryptKey, data[4:])
	return nil
}
