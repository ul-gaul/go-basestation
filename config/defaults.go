package config

import (
    "github.com/spf13/viper"
    "go.bug.st/serial"
    "time"
)

func applyDefaults() {
    viper.SetDefault("Comms.Acknowledge.Timeout", 30*time.Second)
    viper.SetDefault("Comms.Acknowledge.LossThreshold", 20)
    viper.SetDefault("Comms.Acknowledge.BufferSize", 5)
    
    viper.SetDefault("Comms.RocketPacket.LossThreshold", 5)
    viper.SetDefault("Comms.RocketPacket.BufferSize", 128)
    
    viper.SetDefault("Comms.Serial.BaudRate", 115200)
    viper.SetDefault("Comms.Serial.DataBits", 8)
    viper.SetDefault("Comms.Serial.Parity", serial.NoParity)
    viper.SetDefault("Comms.Serial.StopBits", serial.OneStopBit)
    
    viper.SetDefault("Comms.UseBigEndian", false)
    
    viper.SetDefault("Frontend.PlotScale", 1.0)
}
