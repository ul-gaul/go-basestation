package cfg

import (
    "go.bug.st/serial"
    "time"
)

// Comms contains the communication settings
var Comms struct {
    
    // Acknowledge settings
    Acknowledge struct {
        
        // Timeout is the maximum time to wait for an acknowledge after sending
        // a command (see time.ParseDuration)
        Timeout time.Duration
        
        // LossThreshold is the limit of consecutive lost packets before crashing
        LossThreshold uint
        
        // BufferSize is the capacity of the channel for incoming packets
        BufferSize uint
    }
    
    // RocketPacket settings
    RocketPacket struct {
        
        // LossThreshold is the limit of consecutive lost packets before crashing
        LossThreshold uint
        
        // BufferSize is the capacity of the channel for incoming packets
        BufferSize uint
    }
    
    // Serial port settings
    Serial serial.Mode
    
    // UseBigEndian determines the byte order (true: BigEndian, false: LittleEndian).
    UseBigEndian bool
}

// Frontend contains the settings related to the frontend/interface.
var Frontend struct {
    
    ShowFPS bool
    ShowTPS bool
    
    
    // DataRoller contains the settings for the data roller (see vars.DataRoller)
    DataRoller struct {
        // MinTimeGap is the minimum time between packets (see time.ParseDuration)
        MinTimeGap time.Duration
        
        // Limit is the maximum number of packets the DataRoller can contain
        Limit int
    }
    
    // Graph contains the settings for the charts.
    Graph struct {
        
        // BaseWidth is the base width of the charts. Changing this value will
        // change the charts resolution.
        BaseWidth int
        
        // Scale scales the charts.
        Scale float64
    }
}
