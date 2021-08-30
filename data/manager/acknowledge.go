package manager

import (
	"github.com/panjf2000/ants/v2"
	"github.com/ul-gaul/go-basestation/data/packet"
	"github.com/ul-gaul/go-basestation/utils"
	"sync"
)

// TODO documentation

type AcknowledgeCallback func(id uint16, success bool)

var (
	ackLock         sync.Mutex
	ackCallbacks    = make(map[uint16]AcknowledgeCallback)
	unprocessedAcks = make(map[uint16]bool)
)

func SetAcknowledgeCallback(ackId uint16, cb AcknowledgeCallback) {
	ackLock.Lock()
	defer ackLock.Unlock()

	if ack, ok := unprocessedAcks[ackId]; ok {
		utils.CheckErr(ants.Submit(func() {
			cb(ackId, ack)
		}))
		delete(unprocessedAcks, ackId)
	} else {
		ackCallbacks[ackId] = cb
	}
}

func RemoveAcknowledgeCallback(ackId uint16) {
	ackLock.Lock()
	defer ackLock.Unlock()
	delete(ackCallbacks, ackId)
}

func Acknowledge(ack packet.AcknowledgePacket) {
	ackLock.Lock()
	defer ackLock.Unlock()

	success := ack.Result == packet.AckSuccess

	if cb, ok := ackCallbacks[ack.Id]; ok {
		utils.CheckErr(ants.Submit(func() {
			cb(ack.Id, success)
		}))
		delete(ackCallbacks, ack.Id)
	} else {
		unprocessedAcks[ack.Id] = success
	}
}
