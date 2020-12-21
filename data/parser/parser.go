package parser

import (
    "github.com/ul-gaul/go-basestation/data/packet"
)

type ParserFunc func(buffer []byte) (packet.RocketPacket, error)

type ISerialPacketParser interface {
    // Bytes retourne le nombre d'octets d'un paquet
    Bytes() uint
    
    // Parse transforme un buffer en RocketPacket
    Parse(buffer []byte) (packet.RocketPacket, error)
}

var _ ISerialPacketParser = parser{}

type parser struct {
    bytes uint
    parse ParserFunc
}

func Parser(bytes uint, parse ParserFunc) ISerialPacketParser {
    return parser{bytes, parse}
}

func (s parser) Bytes() uint { return s.bytes }

func (s parser) Parse(buffer []byte) (packet.RocketPacket, error) { return s.parse(buffer) }
