package controller

import (
	"encoding/csv"
	"errors"
	"github.com/jszwec/csvutil"
	"github.com/panjf2000/ants/v2"
	log "github.com/sirupsen/logrus"
	"github.com/ul-gaul/go-basestation/data/manager"
	"github.com/ul-gaul/go-basestation/data/packet"
	"github.com/ul-gaul/go-basestation/utils"
	"io"
	"math"
	"time"
)

// TODO documentation

var ErrReplayRunning = errors.New("replay is running")

// Private variables
var (
	replayDecoder                    *csvutil.Decoder
	replayRunning, replayPaused      bool
	replayStartTime, replayPauseTime time.Time
	replayPauseDuration              time.Duration
)

// Public variables
var (
	ReplaySpeed   = 1.0
	OnReplayError func(error)
)

func Replay(reader io.Reader) error {
	var err error
	if IsReplayStarted() {
		return ErrReplayRunning
	}

	if IsConnected() {
		return ErrSerialConnectionOpened
	}

	if replayDecoder, err = csvutil.NewDecoder(csv.NewReader(reader)); err != nil {
		return err
	}

	replayReset()

	replayRunning = true

	chPacket := make(chan packet.RocketPacket, 2)

	utils.CheckErr(ants.Submit(func() {
		replayReadLoop(chPacket)
	}))

	utils.CheckErr(ants.Submit(func() {
		replayDispatchLoop(chPacket)
	}))

	log.Debug("replay started!")
	return nil
}

func IsReplayPaused() bool  { return replayPaused }
func IsReplayStarted() bool { return replayRunning }

func PauseReplay() {
	if replayRunning && !replayPaused {
		replayPauseTime = time.Now()
		replayPaused = true
	}
}

func ResumeReplay() {
	if replayRunning && replayPaused {
		replayPauseDuration += time.Since(replayPauseTime)
		replayPaused = false
	}
}

func ReplayTimeSinceT0() time.Duration {
	if !replayRunning {
		return 0
	}

	d := time.Since(replayStartTime) - replayPauseDuration
	if replayPaused {
		d -= time.Since(replayPauseTime)
	}

	return d
}

func StopReplay() {
	replayRunning = false
	time.Sleep(50 * time.Millisecond)
	replayReset()
}

func replayReset() {
	replayPaused = false
	replayPauseDuration = 0
}

func replayReadLoop(chPacket chan<- packet.RocketPacket) {
	defer close(chPacket)
	var pkt packet.RocketPacket
	pktSent := true

	for replayRunning {
		if pktSent {
			if err := replayDecoder.Decode(&pkt); err == io.EOF {
				log.Debug("simulation ended!")
				return // EOF reached (no more data)
			} else if err != nil {
				if OnReplayError != nil {
					OnReplayError(err)
				}
			} else {
				pktSent = false
			}
		}

		if !pktSent {
			select {
			case chPacket <- pkt:
				pktSent = true

			// Prevent waiting for pkt to be sent so we can keep checking
			// if the simulation has not stopped (running == true)
			default:
			}
		}
	}
}

func replayDispatchLoop(chPacket <-chan packet.RocketPacket) {
	defer func() { replayRunning = false }()

	var pkt packet.RocketPacket
	var ok bool
	pktSent := true

	for replayRunning {
		if !replayPaused && !pktSent {
			t := math.Round(float64(pkt.Time) * float64(time.Millisecond) * ReplaySpeed)
			if ReplayTimeSinceT0() >= time.Duration(t) {
				manager.Data(pkt)
				pktSent = true
			}
		}

		if pktSent {
			select {
			case pkt, ok = <-chPacket:
				if !ok {
					return // chPacket is closed
				}
				pktSent = false

			// Prevent waiting for result from chPacket so we can keep checking
			// if the simulation has not stopped (replayRunning == true)
			default:
			}
		}
	}
}
