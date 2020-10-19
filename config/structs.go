package config

import (
    "go.bug.st/serial"
    "time"
)

var Comms struct {
    AcknowledgeTimeout time.Duration `mapstructure:"acknowledge.timeout"`
    Serial             struct {
        BaudRate        int `mapstructure:"baudrate"`
        DataBits        int `mapstructure:"databits"`
        serial.Parity   `mapstructure:"parity"`
        serial.StopBits `mapstructure:"stopbits"`
    } `mapstructure:"serial"`
}
