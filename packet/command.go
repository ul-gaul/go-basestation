package packet

import (
    "encoding/binary"
    "github.com/ul-gaul/go-basestation/constants"
    "github.com/ul-gaul/go-basestation/utils"
)

type CmdFunction uint8

const (
    FncStartSequence CmdFunction = 0b0001 << iota
    FncSetActuator
    FncResetActuator
    FncSetActuators
    FncResetActuators
    FncSetActuatorsState
)

type Actuator uint8

const (
    ActNone   Actuator = 0b0000
    ActValveA Actuator = 0b0001 << iota
    ActValveB
    ActValveC
    ActValveD
    ActValveE
    ActPiston
)

type CommandPacket struct {
    Id       uint16
    Function CmdFunction
    Argument Actuator
}

func (cp CommandPacket) ToBytes() []byte {
    var buffer []byte
    binary.LittleEndian.PutUint16(buffer, constants.CommandStart)
    binary.LittleEndian.PutUint16(buffer, cp.Id)
    buffer = append(buffer, byte(cp.Function), byte(cp.Argument))
    utils.AppendChecksum(&buffer)
    return buffer
}