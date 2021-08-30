package packet

type CmdFunction uint8

const (
    FncStartSequence CmdFunction = 0b0000_0001 << iota
    FncSetActuator
    FncResetActuator
    FncSetActuators
    FncResetActuators
    FncSetActuatorsState
)

type Actuator uint8

const (
    ActValveA Actuator = 0b0000_0001 << iota
    ActValveB
    ActValveC
    ActValveD
    ActValveE
    ActPiston
    ActNone Actuator = 0b0000_0000
)

type CommandPacket struct {
    Id       uint16
    Function CmdFunction
    Argument Actuator
}
