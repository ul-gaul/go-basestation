package internal

import (
	"bytes"
	"encoding/binary"
	e "errors"
	"fmt"
	"github.com/sigurn/crc16"
	"github.com/ul-gaul/go-basestation/cfg"
	"github.com/ul-gaul/go-basestation/data/packet"
	"io"
)

const (
	StartDelimSize = 2
	ChecksumSize   = 2

	RocketPacketStart uint16 = 's'
	AcknowledgeStart  uint16 = 0xface
	CommandStart      uint16 = 0xface
)

var ErrInvalidChecksum = e.New("packet checksum validation failed")

var crcTable = crc16.MakeTable(crc16.CRC16_CCITT_FALSE)

var _ error = ErrUnknownStart{}

type ErrUnknownStart struct {
	Char uint16
}

func (e ErrUnknownStart) Error() string {
	return fmt.Sprintf("unknown start: %x", e.Char)
}

// ReadBinary reads and unmarshall the next packet from a raw binary io.Reader
func ReadBinary(reader io.Reader) (result interface{}, err error) {
	order := binary.ByteOrder(binary.LittleEndian)
	if cfg.Comms.UseBigEndian {
		order = binary.BigEndian
	}

	// Read the packet start delimiter
	start := make([]byte, StartDelimSize)
	if _, err = reader.Read(start); err != nil {
		return nil, err
	}

	// Determine which type of packet we are reading
	switch char := order.Uint16(start); char {
	case RocketPacketStart:
		result = packet.RocketPacket{}
	case AcknowledgeStart:
		result = packet.AcknowledgePacket{}
	default:
		return nil, ErrUnknownStart{char}
	}

	// Read packet data
	bufData := make([]byte, binary.Size(result))
	if _, err = reader.Read(bufData); err != nil {
		return nil, err
	}

	// Read checksum
	bufChecksum := make([]byte, ChecksumSize)
	if _, err = reader.Read(bufChecksum); err != nil {
		return nil, err
	}

	// Calculate checksum of data and compare result with checksum from reader
	if crc16.Checksum(bufData, crcTable) != order.Uint16(bufChecksum) {
		return nil, ErrInvalidChecksum
	}

	// Unmarshall data
	err = binary.Read(bytes.NewReader(bufData), order, &result)

	return result, err
}

// WriteCommandBinary TODO doc
func WriteCommandBinary(writer io.Writer, cmdPacket packet.CommandPacket) error {
	return WriteBinary(writer, CommandStart, cmdPacket)
}

// WriteBinary TODO doc
func WriteBinary(writer io.Writer, delimiter uint16, packet interface{}) (err error) {
	order := binary.ByteOrder(binary.LittleEndian)
	if cfg.Comms.UseBigEndian {
		order = binary.BigEndian
	}

	buffer := new(bytes.Buffer)

	// Write the delimiter
	if err = binary.Write(buffer, order, delimiter); err != nil {
		return err
	}

	// Write the data/packet
	if err = binary.Write(buffer, order, packet); err != nil {
		return err
	}

	// Calculate and write the checksum
	checksum := crc16.Checksum(buffer.Bytes()[StartDelimSize:], crcTable)
	if err = binary.Write(buffer, order, checksum); err != nil {
		return err
	}

	_, err = writer.Write(buffer.Bytes())
	return err
}
