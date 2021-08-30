package packet

type AcknowledgeResult uint8

const (
	AckSuccess AcknowledgeResult = 0x01
	AckFailure AcknowledgeResult = 0xFF
)

type AcknowledgePacket struct {
	Id     uint16
	Result AcknowledgeResult
}
