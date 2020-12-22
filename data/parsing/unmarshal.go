package parsing

import (
    "bytes"
    "encoding/binary"
    "errors"
    
    "github.com/ul-gaul/go-basestation/config"
    "github.com/ul-gaul/go-basestation/data/packet"
)

func UnmarshalParser() (ISerialPacketParser, error) {
    order := binary.ByteOrder(binary.LittleEndian)
    if config.Comms.UseBigEndian {
        order = binary.BigEndian
    }
    
    size := binary.Size(packet.RocketPacket{})
    if size < 0 {
        return nil, errors.New("invalid byte size of RocketPacket")
    }
    return Parser(uint(size), func(buffer []byte) (packet.RocketPacket, error) {
        var pkt packet.RocketPacket
        err := binary.Read(bytes.NewReader(buffer), order, &pkt)
        return pkt, err
    }), nil
}
