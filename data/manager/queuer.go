package manager

import (
	"github.com/panjf2000/ants/v2"
	"github.com/ul-gaul/go-basestation/data/packet"
	"sync"
)

type queuer struct {
	data    []packet.RocketPacket
	handler IDataHandler
	running bool
	sync.Mutex
}

func (q *queuer) tryHandle() {
	var data []packet.RocketPacket
	q.Lock()
	if q.running || len(q.data) == 0 {
		q.Unlock()
		return
	}
	q.running = true
	data = q.data
	q.data = []packet.RocketPacket{}
	q.Unlock()

	defer func() {
		q.Lock()
		q.running = false
		q.Unlock()
		_ = ants.Submit(q.tryHandle)
	}()

	q.handler.OnData(data)
}

func (q *queuer) Add(data []packet.RocketPacket) {
	q.Lock()
	q.data = append(q.data, data...)
	q.Unlock()
	defer ants.Submit(q.tryHandle)
}
