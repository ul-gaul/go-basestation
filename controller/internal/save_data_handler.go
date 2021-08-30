package internal

import (
	"github.com/jszwec/csvutil"
	log "github.com/sirupsen/logrus"
	"github.com/ul-gaul/go-basestation/data/manager"
	"github.com/ul-gaul/go-basestation/data/packet"
)

var _ manager.IDataHandler = (*SaveDataHandler)(nil)

type SaveDataHandler struct {
	Output  *csvutil.Encoder
	Enabled bool
}

func (s *SaveDataHandler) OnData(packets []packet.RocketPacket) {
	if s.Enabled && s.Output != nil {
		if err := s.Output.Encode(packets); err != nil {
			log.Error(err)
		}
	}
}
