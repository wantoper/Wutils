package Msg

import (
	"errors"
	"unsafe"
)

type ClientHello struct {
	Version      uint16
	SuiteLength  uint8
	CipherSuites []uint8
}

func (c *ClientHello) Marshal() []byte {
	length := int(unsafe.Sizeof(c.Version)) + int(unsafe.Sizeof(c.SuiteLength)) + len(c.CipherSuites)
	data := make([]byte, length)
	data[0] = byte(c.Version >> 8)
	data[1] = byte(c.Version & 0xFF)
	data[2] = c.SuiteLength
	copy(data[3:], c.CipherSuites)

	return data
}

func (c *ClientHello) UnmarShal(data []byte) error {
	if len(data) < 4 {
		return errors.New("too short data")
	}

	c.Version = uint16(data[0])<<8 | uint16(data[1])
	c.SuiteLength = data[2]

	c.CipherSuites = make([]uint8, c.SuiteLength)
	copy(c.CipherSuites, data[3:])
	return nil
}
