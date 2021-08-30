package internal

import (
	"github.com/jszwec/csvutil"
	"io"
	"os"

	"github.com/ul-gaul/go-basestation/data/packet"
)

// ReadCsvFile reads a CSV file and returns a list of packet.RocketPacket.
//
// IMPORTANT: Do NOT use on large files. Instead use csvutil.Decoder and
// process the file line by line (don't store the whole file in memory).
func ReadCsvFile(csvPath string) ([]packet.RocketPacket, error) {
	file, err := os.OpenFile(csvPath, os.O_RDONLY, 0o755)
	if err != nil {
		return nil, err
	}
	return ReadCsv(file)
}

// ReadCsv reads and transforms all the data from the reader to a list
// of packet.RocketPacket.
//
// IMPORTANT: Do NOT use on large files. Instead use csvutil.Decoder and
// process the file line by line (don't store the whole file in memory).
func ReadCsv(reader io.Reader) ([]packet.RocketPacket, error) {
	// Read all bytes
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	// Transform the bytes/data to a list of packets
	var packets []packet.RocketPacket
	err = csvutil.Unmarshal(data, &packets)

	return packets, err
}
