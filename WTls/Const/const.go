package Const

const (
	Version1_0 uint16 = 0x01
	Version1_1 uint16 = 0x02
)

const (
	AES_GCM_128 uint8 = 0x01
	AES_GCM_256 uint8 = 0x02
)

func GetCipherSuiteName(suite uint8) string {
	switch suite {
	case AES_GCM_128:
		return "AES-GCM-128"
	case AES_GCM_256:
		return "AES-GCM-256"
	default:
		return "Unknown Cipher Suite"
	}
}
