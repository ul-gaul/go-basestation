package constants

import (
    "errors"
)

var (
    ErrCommunicatorAlreadyStarted = errors.New("the serial communicator has already been started")
    ErrInvalidChecksum = errors.New("invalid packet checksum")
    ErrAcknowledgeTimeout = errors.New("") // FIXME
    ErrAcknowledgeFail = errors.New("") // FIXME
    ErrLostTooManyRocketPacket = errors.New("lost too many RocketPacket packets")
    ErrLostTooManyAcknowledge = errors.New("lost too many Acknowledge packets")
    
    ErrNotARegularFile = errors.New("file must be a regular file")
)


/** UI Errors ***/
var (
    ErrPaddingOutOfRange = errors.New("padding ratio must be between -1 and 1")
)