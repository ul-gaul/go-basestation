package controller

// TODO documentation

type Mode uint8

const (
    MODE_SERIAL Mode = 0b0000_0001 << iota
    MODE_REPLAY
    MODE_GENERATE
    MODE_NONE Mode = 0b0000_0000
)

func CurrentMode() Mode {
    if IsConnected() {
        return MODE_SERIAL
    }
    if IsReplayStarted() {
        return MODE_REPLAY
    }
    if IsGeneratorStarted() {
        return MODE_GENERATE
    }
    return MODE_NONE
}
