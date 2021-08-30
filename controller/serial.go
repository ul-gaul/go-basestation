package controller

import (
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/jszwec/csvutil"
	"github.com/panjf2000/ants/v2"
	log "github.com/sirupsen/logrus"
	"github.com/ul-gaul/go-basestation/cfg"
	"github.com/ul-gaul/go-basestation/controller/internal"
	"github.com/ul-gaul/go-basestation/data/manager"
	"github.com/ul-gaul/go-basestation/data/packet"
	"github.com/ul-gaul/go-basestation/utils"
	"go.bug.st/serial"
	"os"
	"path/filepath"
	"time"
)

const DefaultOutputDir = "./output"

var ErrSerialConnectionOpened = errors.New("serial connection opened")
var ErrOutputChangedWhileRunning = errors.New("cannot change output file while running")

var (
	OnSerialError       func(error)
	serialOutputFile    *os.File
	conn                serial.Port
	serialOutputHandler = &internal.SaveDataHandler{}
)

func init() {
	var err error
	utils.CheckErr(os.MkdirAll(DefaultOutputDir, 0755))
	path := filepath.Join(DefaultOutputDir, fmt.Sprintf("serial_%s.csv", time.Now().Format("2006-01-02T15.04.05")))
	serialOutputFile, err = os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0755)
	utils.CheckErr(err)

	manager.AddDataHandler(serialOutputHandler)
}

// IsConnected returns whether a serial connection is open or not.
func IsConnected() bool { return conn != nil }

// CloseConnection closes the connection to the serial port.
// It will NOT panic if there's no opened connection.
func CloseConnection() {
	if conn != nil {
		if err := conn.Close(); err != nil {
			log.Warn(err)
		}
		conn = nil
	}
}

// ListAvailablePorts lists all available serial ports.
//
// See serial.GetPortsList
func ListAvailablePorts() []string {
	ports, err := serial.GetPortsList()
	utils.CheckErr(err)
	return ports
}

// OpenConnection opens a connection to the specified port and starts listenning
// for incoming data.
func OpenConnection(port string) (err error) {
	if IsReplayStarted() {
		return ErrReplayRunning
	}

	if IsConnected() {
		return ErrSerialConnectionOpened
	}

	conn, err = serial.Open(port, &cfg.Comms.Serial)
	if err != nil {
		conn = nil
		return err
	}

	serialOutputHandler.Output = csvutil.NewEncoder(csv.NewWriter(serialOutputFile))
	serialOutputHandler.Enabled = true

	utils.CheckErr(ants.Submit(listen))

	return nil
}

func listen() {
	defer func() {
		_ = conn.Close()
		conn = nil
		serialOutputHandler.Enabled = false
	}()

	for {
		pkt, err := internal.ReadBinary(conn)
		if err == os.ErrClosed {
			return
		} else if err != nil {
			OnSerialError(err)
		}

		switch pkt := pkt.(type) {
		case packet.RocketPacket:
			manager.Data(pkt)
		case packet.AcknowledgePacket:
			manager.Acknowledge(pkt)
		default:
			log.Warnf("Unkown packet type: %s\n", pkt)
		}
	}
}

func SetSerialOutputFile(path string, overwrite bool) error {
	if path == serialOutputFile.Name() {
		return nil
	}

	if serialOutputHandler.Enabled {
		return ErrOutputChangedWhileRunning
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	flag := os.O_CREATE | os.O_WRONLY
	if overwrite {
		flag |= os.O_TRUNC
	} else {
		flag |= os.O_EXCL
	}

	file, err := os.OpenFile(path, flag, 0755)
	if err != nil {
		return err
	}

	stat, err := file.Stat()
	if err != nil {
		return err
	}

	if stat.IsDir() {
		return ErrNotARegularFile
	}

	serialOutputFile = file
	return nil
}

func GetSerialOutputFile() *os.File {
	return serialOutputFile
}