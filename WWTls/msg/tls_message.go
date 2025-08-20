package msg

type TlsVersion uint16

const (
	TlsVersion_1_0 TlsVersion = 0x0301
	TlsVersion_1_1 TlsVersion = 0x0302
	TlsVersion_1_2 TlsVersion = 0x0303
	TlsVersion_1_3 TlsVersion = 0x0304
)

type ContentType uint8

const (
	RecordMessage_ContentType_Handshake ContentType = 22
)

// TlsRecord层信息
// TLSCiphertext的实现: HandshakeMessage
type RecordMessage struct {
	ContentType ContentType
	Version     TlsVersion
	Length      uint16
	//TLSCiphertext []byte
}

type HandshakeType uint8

const (
	HandshakeMessage_HandshakeType_ClientHello HandshakeType = 1
	HandshakeMessage_HandshakeType_ServerHello HandshakeType = 2
	HandshakeMessage_HandshakeType_Certificate HandshakeType = 11
	HandshakeMessage_HandshakeType_Finished    HandshakeType = 20
)

// Handshake信息
type HandshakeMessage struct {
	HandshakeType HandshakeType
	Length        uint32
	//HandshakeMessage []byte : HandshakeMessageI.marshal()
}

type HandshakeMessageI interface {
	marshal() ([]byte, error)
	unmarshal([]byte) bool
}

type HandshakeMessageClientHello struct {
	LegacyVersion TlsVersion
	Random        [32]byte
}
