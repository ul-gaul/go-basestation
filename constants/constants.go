package constants

import (
    log "github.com/sirupsen/logrus"
)

const (
    StartDelimSize = 2
    AckPacketSize  = 5 // Avec CRC
    
    RocketPacketStart uint16 = 's'
    AcknowledgeStart  uint16 = 0xface
    CommandStart      uint16 = 0xface
    
    DefaultLogLevel = log.InfoLevel
)
