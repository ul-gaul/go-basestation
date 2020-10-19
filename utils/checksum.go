package utils

import (
    "encoding/binary"
    "github.com/sigurn/crc16"
)
var table = crc16.MakeTable(crc16.CRC16_CCITT_FALSE)

// ValidateChecksum valide le checksum des données spécifiées.
// Les 2 derniers octets des données spéciées doit contenir le CRC.
func ValidateChecksum(data []byte) bool {
    crc := binary.LittleEndian.Uint16(data[len(data)-2:])
    return crc16.Checksum(data[:len(data)-2], table) == crc
}

// AppendChecksum calcule le CRC des données spécifiées et l'ajoute à celles-ci.
func AppendChecksum(data *[]byte) {
    crc := crc16.Checksum(*data, table)
    rawCrc := make([]byte, 2)
    binary.LittleEndian.PutUint16(rawCrc, crc)
    *data = append(*data, rawCrc...)
}