package persistence

import (
    "encoding/binary"
    "errors"
    "github.com/panjf2000/ants/v2"
    log "github.com/sirupsen/logrus"
    "go.bug.st/serial"
    "sync"
    "time"
    
    . "github.com/ul-gaul/go-basestation/config"
    "github.com/ul-gaul/go-basestation/constants"
    "github.com/ul-gaul/go-basestation/data/packet"
    . "github.com/ul-gaul/go-basestation/data/parser"
    "github.com/ul-gaul/go-basestation/utils"
)

type ISerialPacketCommunicator interface {
    serial.Port
    
    // Start ouvre le port et démarre le thread qui lit les données reçues
    Start() error
    
    // SendCommand envoi une commande à la fusé et attend le acknowledge.
    // L'erreur ErrAcknowledgeFail est retournée lorsque le acknowledge échou.
    SendCommand(fnc packet.CmdFunction, arg packet.Actuator) error
    
    // RocketPacketChannel retourne le channel auquel sont envoyés les RocketPaquets reçus de la fusé.
    RocketPacketChannel() <-chan packet.RocketPacket
    
    // ErrorChannel retourne le channel auquel sont envoyés les erreurs qui surviennent dans le thread de lecture.
    ErrorChannel() <-chan error
}

type serialPacketCommunicator struct {
    parser  ISerialPacketParser
    strPort string
    lastCmdId uint16
    mut       sync.Mutex
    serial.Port
    
    chAcknowledge  chan packet.AcknowledgePacket
    chRocketPacket chan packet.RocketPacket
    chError        chan error
}

func NewSerialPacketCommunicator(port string, parser ISerialPacketParser) ISerialPacketCommunicator {
    return &serialPacketCommunicator{
        parser:         parser,
        strPort:        port,
        chAcknowledge:  make(chan packet.AcknowledgePacket, Comms.Acknowledge.BufferSize),
        chRocketPacket: make(chan packet.RocketPacket, Comms.RocketPacket.BufferSize),
        chError:        make(chan error),
    }
}

func (s *serialPacketCommunicator) Start() (err error) {
    if s.Port != nil {
        err = constants.ErrCommunicatorAlreadyStarted
    } else if s.Port, err = serial.Open(s.strPort, &Comms.Serial); err == nil {
        err = ants.Submit(s.run)
        if err != nil {
            log.Fatal(err)
        }
    }
    return err
}

func (s *serialPacketCommunicator) run() {
    defer s.Port.Close()
    var err error
    var responseType uint16
    var pkt packet.RocketPacket
    var ack packet.AcknowledgePacket
    var pktLost, ackLost uint
    start := make([]byte, constants.StartDelimSize)
    
    for ; err == nil; _, err = s.Port.Read(start) {
        responseType = binary.LittleEndian.Uint16(start)
        
        switch responseType {
        case constants.RocketPacketStart:
            if pkt, err = s.readRocketPacket(); err == nil {
                select {
                case s.chRocketPacket <- pkt:
                    pktLost = 0
                default:
                    pktLost++
                    if pktLost > Comms.RocketPacket.LossThreshold {
                        log.Panicln(constants.ErrLostTooManyRocketPacket)
                    } else {
                        log.Warnln("Packet lost -> RocketPacket")
                    }
                }
            }
        case constants.AcknowledgeStart:
            if ack, err = s.readAcknowledge(); err == nil {
                select {
                case s.chAcknowledge <- ack:
                    ackLost = 0
                default:
                    ackLost++
                    if ackLost > Comms.Acknowledge.LossThreshold {
                        log.Panicln(constants.ErrLostTooManyAcknowledge)
                    } else {
                        log.Warnln("Packet lost -> Acknowledge")
                    }
                }
            }
        default:
            log.Warnf("Unknown start: %x\n", responseType)
        }
    }
    
    var portErr serial.PortError
    if errors.As(err, &portErr) && portErr.Code() != serial.PortClosed {
        s.chError <- err
    }
}

func (s *serialPacketCommunicator) readRocketPacket() (packet.RocketPacket, error) {
    var pkt packet.RocketPacket
    var err error
    
    buffer := make([]byte, s.parser.Bytes()+2)
    if _, err = s.Read(buffer); err == nil {
        if !utils.ValidateChecksum(buffer) {
            err = constants.ErrInvalidChecksum
        } else {
            pkt, err = s.parser.Parse(buffer[:len(buffer)-2])
        }
    }
    
    return pkt, err
}

func (s *serialPacketCommunicator) readAcknowledge() (packet.AcknowledgePacket, error) {
    var ack packet.AcknowledgePacket
    var err error
    
    buffer := make([]byte, constants.AckPacketSize)
    if _, err = s.Read(buffer); err == nil {
        if !utils.ValidateChecksum(buffer) {
            err = constants.ErrInvalidChecksum
        } else {
            buffer = buffer[:len(buffer)-2]
            ack = packet.AcknowledgePacket{
                Id:  binary.LittleEndian.Uint16(buffer),
                Ack: packet.Acknowledge(buffer[2]),
            }
        }
    }
    
    return ack, err
}

func (s *serialPacketCommunicator) SendCommand(fnc packet.CmdFunction, arg packet.Actuator) error {
    var err error
    var ack packet.AcknowledgePacket
    
    s.mut.Lock()
    s.lastCmdId++
    id := s.lastCmdId
    s.mut.Unlock()
    
    _, err = s.Write(packet.CommandPacket{Id: id, Function: fnc, Argument: arg}.ToBytes())
    
    for timeout := time.After(Comms.Acknowledge.Timeout); ack.Id != id && err == nil; {
        select {
        case ack = <-s.chAcknowledge:
            if ack.Ack != packet.AckSuccess {
                err = constants.ErrAcknowledgeFail
            }
        case <-timeout:
            err = constants.ErrAcknowledgeTimeout
        default:
        }
    }
    
    return err
}

func (s *serialPacketCommunicator) RocketPacketChannel() <-chan packet.RocketPacket {
    return s.chRocketPacket
}

func (s *serialPacketCommunicator) ErrorChannel() <-chan error {
    return s.chError
}
