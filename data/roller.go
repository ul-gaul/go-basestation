package data

import (
	"container/ring"
	"github.com/ul-gaul/go-basestation/data/manager"
	"github.com/ul-gaul/go-basestation/data/packet"
	"math"
	"time"
)

// TODO doc

const (
	DefaultMinTimeGap = 50 * time.Millisecond
	DefaultRollerLimit = 100
)

var _ manager.IDataHandler = (*Roller)(nil)

type Roller struct {
	ring       *ring.Ring
	len        int
	MinTimeGap time.Duration
	lastT      time.Duration
}

func NewRoller(minTimeGap time.Duration, limit int) *Roller {
	if minTimeGap == 0 {
		minTimeGap = DefaultMinTimeGap
	}
	
	if limit == 0 {
		limit = DefaultRollerLimit
	}
	
	return &Roller{ring: ring.New(limit), MinTimeGap: minTimeGap}
}

func (r *Roller) Limit() int {
	return r.ring.Len()
}

func (r *Roller) SetLimit(limit int) {
	limit = int(math.Max(0, float64(limit)))
	if limit > r.ring.Len() {
		r.ring.Link(ring.New(limit - r.ring.Len()))
	} else {
		r.ring.Unlink(r.ring.Len() - limit)
	}

	if limit < r.len {
		r.len = limit
	}
}

func (r *Roller) OnData(packets []packet.RocketPacket) {
	var added int
	empty := r.len == 0

	for _, pkt := range packets {
		if empty {
			empty = false
			r.ring.Value = pkt
			r.ring = r.ring.Next()
			r.ring.Value = pkt
			r.ring = r.ring.Next()
			r.lastT = pkt.Time.Duration()
			added += 2
		} else {
			t := pkt.Time.Duration()
			r.ring.Prev().Value = pkt
			if t-r.lastT >= r.MinTimeGap {
				r.ring.Value = pkt
				r.ring = r.ring.Next()
				r.lastT = t
				added++
			}
		}
	}

	i := r.len + added
	if i > r.Limit() {
		i = r.Limit()
	}
	r.len = i
}

func (r *Roller) Clear() {
	current := r.ring
	for i := 0; i < r.ring.Len(); i++ {
		current.Value = nil
		current = current.Next()
	}
	r.len = 0
}

func (r *Roller) Last() packet.RocketPacket {
	if r.Len() == 0 {
		return packet.RocketPacket{}
	}

	cur := r.ring.Prev()
	for cur.Value == nil && cur != r.ring {
		cur = cur.Prev()
	}
	return cur.Value.(packet.RocketPacket)
}

func (r *Roller) Do(fnc func(packet.RocketPacket)) {
	r.ring.Do(func(i interface{}) {
		if i != nil {
			fnc(i.(packet.RocketPacket))
		}
	})
}

func (r *Roller) Len() int {
	return r.len
}