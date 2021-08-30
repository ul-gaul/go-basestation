package controller

import (
	"errors"
	"github.com/ul-gaul/go-basestation/controller/internal"
	"github.com/ul-gaul/go-basestation/data/manager"
	"github.com/ul-gaul/go-basestation/data/packet/state"
	"os"
)

// TODO documentation

var ErrNotARegularFile = errors.New("file must be a regular file")

// LoadFile TODO doc
func LoadFile(csvFile string) error {
	stat, err := os.Stat(csvFile)
	if err != nil {
		return err
	}

	if !stat.Mode().IsRegular() {
		return ErrNotARegularFile
	}

	packets, err := internal.ReadCsvFile(csvFile)
	if err != nil {
		return err
	}
	manager.SetStaticData(packets)
	return nil
}

func Shutdown() {
	defer CloseConnection()
	defer StopReplay()
	defer StopGenerator()
	defer serialOutputFile.Close()
}


func SetActuator(item state.Item, val bool, callback func()) {
	switch CurrentMode() {
	case MODE_GENERATE:
		rdmState = rdmState.WithStatusOf(item, val)
		callback()
	case MODE_SERIAL:
		fallthrough // TODO send command
	default:
		callback()
	}
}
