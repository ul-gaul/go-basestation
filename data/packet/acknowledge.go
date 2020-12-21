package packet

type Acknowledge uint8

const (
    AckSuccess Acknowledge = 0x01
    AckFailure Acknowledge = 0xFF
)

type AcknowledgePacket struct {
    Id  uint16
    Ack Acknowledge
}
