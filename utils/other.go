package utils

import (
	log "github.com/sirupsen/logrus"
	"math"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
)

// CheckErr logs the error and panics if err is not nil.
func CheckErr(err error) {
	if err != nil {
		if log.IsLevelEnabled(log.DebugLevel) {
			log.Panicf("%s\n\t%s\n", err,
				strings.ReplaceAll(string(debug.Stack()), "\n", "\n\t"))
		} else {
			log.Panic(err)
		}
	}
}

// GetExecutablePath TODO doc
func GetExecutablePath() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		exe, err = filepath.EvalSymlinks(exe)
		if err != nil {
			return exe, err
		}
	}
	return filepath.Dir(exe), nil
}

// Constrain ensures value is within the range [min, max].
func Constrain(value, min, max float64) float64 {
	return math.Max(min, math.Min(max, value))
}

// ConstrainInt ensures value is within the range [min, max].
func ConstrainInt(value, min, max int) int {
	return int(Constrain(float64(value), float64(min), float64(max)))
}
