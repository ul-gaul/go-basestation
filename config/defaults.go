package config

import (
    "github.com/spf13/viper"
    "go.bug.st/serial"
    "time"
)

func applyDefaults() {
    viper.SetDefault("comms.acknowledge.timeout", 30 * time.Second)
    viper.SetDefault("comms.serial.baudrate", 115200)
    viper.SetDefault("comms.serial.databits", 8)
    viper.SetDefault("comms.serial.parity", serial.NoParity)
    viper.SetDefault("comms.serial.stopbits", serial.OneStopBit)
}
