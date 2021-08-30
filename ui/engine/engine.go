package engine

import (
	"github.com/ul-gaul/go-basestation/ui/vars"
	"os"
	"time"
)

const DefaultMaxTPS = 25

var (
	running, closed bool

	tps, fps int

	tickMinDeltaT time.Duration
)

// FPS returns the number of frames per second.
func FPS() int { return fps }

// TPS returns the number of ticks per second.
func TPS() int { return tps }

// GetMaxTPS returns the maximum number of ticks per second.
//
// A value of 0 means unlimited.
func GetMaxTPS() int {
	if tickMinDeltaT > 0 {
		return int(time.Second / tickMinDeltaT)
	}
	return 0
}

// SetMaxTPS sets the maximum number of ticks per second.
//
// A negative value or a value of 0 means unlimited.
func SetMaxTPS(tps int) {
	if tps <= 0 {
		tickMinDeltaT = 0
	} else {
		tickMinDeltaT = time.Second / time.Duration(tps)
	}
}

func Close() {
	closed = true
	defer vars.Window.Close()
	defer os.Exit(0)
}

func init() {
	SetMaxTPS(DefaultMaxTPS)
}
