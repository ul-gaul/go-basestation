package constants

const (
    StartDelimSize = 2
    AckPacketSize  = 5 // Avec CRC
    
    RocketPacketStart uint16 = 's'
    AcknowledgeStart  uint16 = 0xface
    CommandStart      uint16 = 0xface
)