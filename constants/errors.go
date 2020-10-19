package constants

import (
    "errors"
)

var (
    ErrCommunicatorAlreadyStarted = errors.New("the serial communicator has already been started")
    ErrInvalidChecksum = errors.New("invalid packet checksum")
    ErrAcknowledgeTimeout = errors.New("") // FIXME
    ErrAcknowledgeFail = errors.New("") // FIXME
)
