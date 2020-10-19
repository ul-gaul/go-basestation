package utils

import (
    "encoding/binary"
    "github.com/stretchr/testify/assert"
    "testing"
)

type TestData struct {
    data []byte
    expectedCrc uint16
}

var tests = []TestData{
    // CRCs obtenu sur https://crccalc.com/
    { []byte("abcdef_xyz"), 0xB918 },
    { []byte{ 0x00 }, 0xE1F0 },
    { []byte{ 0x00, 0x00, 0x00, 0x00 }, 0x84C0 },
}

func TestAppendChecksum(t *testing.T) {
    assert := assert.New(t)
    
    for _, tt := range tests {
        crcBytes := make([]byte, 2)
        binary.LittleEndian.PutUint16(crcBytes, tt.expectedCrc)
        
        expectedData := append(tt.data, crcBytes...)
        testedData := tt.data[:]
        AppendChecksum(&testedData)
        
        assert.Len(testedData, len(expectedData))
        assert.Equal(expectedData, testedData)
    }
}

func TestValidateChecksum(t *testing.T) {
    assert := assert.New(t)
    
    for _, tt := range tests {
        data := append(tt.data, 0, 0)
        
        binary.LittleEndian.PutUint16(data[len(data)-2:], tt.expectedCrc)
        assert.True(ValidateChecksum(data))
        
        binary.LittleEndian.PutUint16(data[len(data)-2:], tt.expectedCrc + 1)
        assert.False(ValidateChecksum(data))
    
        binary.LittleEndian.PutUint16(data[len(data)-2:], tt.expectedCrc - 1)
        assert.False(ValidateChecksum(data))
    }
}