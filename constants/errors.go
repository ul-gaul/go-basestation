package constants

import (
    e "errors"
)

var (
    ErrCommunicatorAlreadyStarted = e.New("the serial communicator has already been started")
    ErrInvalidChecksum = e.New("invalid packet checksum")
    ErrAcknowledgeTimeout = e.New("") // FIXME
    ErrAcknowledgeFail = e.New("")    // FIXME
    ErrLostTooManyRocketPacket = e.New("lost too many RocketPacket packets")
    ErrLostTooManyAcknowledge = e.New("lost too many Acknowledge packets")
    
    ErrNotARegularFile = e.New("file must be a regular file")
    ErrConnectionAlreadyOpen = e.New("another connection is open")
)


/** UI Errors ***/
var (
    ErrPaddingOutOfRange = e.New("padding ratio must be between -1 and 1")
)