package config

import (
    "go.bug.st/serial"
    "time"
)

var Comms struct {
    
    Acknowledge struct {
        Timeout time.Duration
        LossThreshold uint
        BufferSize uint
    }
    
    RocketPacket struct {
        LossThreshold uint
        BufferSize uint
    }
    
    Serial serial.Mode
    UseBigEndian bool
}
