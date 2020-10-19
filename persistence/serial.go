package persistence

import (
    "encoding/binary"
    "errors"
    "github.com/panjf2000/ants/v2"
    . "github.com/ul-gaul/go-basestation/config"
    "github.com/ul-gaul/go-basestation/constants"
    "github.com/ul-gaul/go-basestation/packet"
    "github.com/ul-gaul/go-basestation/utils"
    "go.bug.st/serial"
    "log"
    "sync"
    "time"
)

type ISerialPacketParser interface {
    // Bytes retourne le nombre d'octets d'un paquet
    Bytes() uint
    
    // Parse transforme un buffer en RocketPacket
    Parse(buffer []byte) (packet.RocketPacket, error)
}

type ISerialPacketCommunicator interface {
    serial.Port
    
    // Start ouvre le port et démarre le thread qui lit les données reçues
    Start() error
    
    // SendCommand envoi une commande à la fusé et attend le acknowledge.
    // L'erreur ErrAcknowledgeFail est retournée lorsque le acknowledge échou.
    SendCommand(fnc packet.CmdFunction, arg packet.Actuator) error
    
    // RocketPacketChannel retourne le channel auquel sont envoyés les RocketPaquets reçus de la fusé.
    RocketPacketChannel() <- chan packet.RocketPacket
    
    // ErrorChannel retourne le channel auquel sont envoyés les erreurs qui surviennent dans le thread de lecture.
    ErrorChannel() <- chan error
}

type serialPacketCommunicator struct {
    parser ISerialPacketParser
    strPort string
    lastCmdId uint16
    mut sync.Mutex
    serial.Port
    
    chAcknowledge  chan packet.AcknowledgePacket
    chRocketPacket chan packet.RocketPacket
    chError        chan error
}

func NewSerialPacketCommunicator(port string, parser ISerialPacketParser) ISerialPacketCommunicator {
    return &serialPacketCommunicator{
        parser:  parser,
        strPort: port,
        chAcknowledge: make(chan packet.AcknowledgePacket),
        chRocketPacket: make(chan packet.RocketPacket),
        chError: make(chan error),
    }
}

func (s *serialPacketCommunicator) Start() error {
    var err error
    mode := serial.Mode{
        Comms.Serial.BaudRate,
        Comms.Serial.DataBits,
        Comms.Serial.Parity,
        Comms.Serial.StopBits,
    }
    
    if s.Port != nil {
        err = constants.ErrCommunicatorAlreadyStarted
    } else if s.Port, err = serial.Open(s.strPort, &mode); err == nil {
        err = ants.Submit(s.run)
        if err != nil {
            log.Panicln(err)
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
    start := make([]byte, constants.StartDelimSize)
    
    for ; err == nil; _, err = s.Port.Read(start) {
        responseType = binary.LittleEndian.Uint16(start)
        
        switch responseType {
        case constants.RocketPacketStart:
            if pkt, err = s.readRocketPacket(); err == nil {
                s.chRocketPacket <- pkt
            }
        case constants.AcknowledgeStart:
            if ack, err = s.readAcknowledge(); err == nil {
                s.chAcknowledge <- ack
            }
        default:
            log.Printf("Unknown start: %x\n", responseType)
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
    
    _, err = s.Write(packet.CommandPacket{Id: id, Function: fnc, Argument: arg }.ToBytes())
    
    for timeout := time.After(Comms.AcknowledgeTimeout); ack.Id != id && err == nil; {
        select {
        case ack = <- s.chAcknowledge:
            if ack.Ack != packet.AckSuccess {
                err = constants.ErrAcknowledgeFail
            }
        case <- timeout:
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
