package manager

import (
	"github.com/ul-gaul/go-basestation/data/packet"
	"sync"
)

// TODO documentation

type IDataHandler interface {
	OnData([]packet.RocketPacket)
}

var (
	dataLock sync.RWMutex
	queuers  = make(map[IDataHandler]*queuer)

	static []packet.RocketPacket
)

func AddDataHandler(handlers ...IDataHandler) {
	dataLock.Lock()
	defer dataLock.Unlock()

	for _, handler := range handlers {
		if _, exists := queuers[handler]; !exists {
			queuers[handler] = &queuer{handler: handler}
		}
	}
}

func RemoveDataHandler(handlers ...IDataHandler) {
	dataLock.Lock()
	defer dataLock.Unlock()

	for _, handler := range handlers {
		delete(queuers, handler)
	}
}

func Data(data ...packet.RocketPacket) {
	dataLock.RLock()
	defer dataLock.RUnlock()

	for _, queuer := range queuers {
		queuer.Add(data)
	}
}

func SetStaticData(packets []packet.RocketPacket) {
	dataLock.Lock()
	defer dataLock.Unlock()
	static = packets
}

func GetStaticData() []packet.RocketPacket {
	dataLock.RLock()
	defer dataLock.RUnlock()
	return static
}
