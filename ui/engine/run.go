package engine

import (
	"gioui.org/app"
	log "github.com/sirupsen/logrus"
	"github.com/ul-gaul/go-basestation/pool"
	"github.com/ul-gaul/go-basestation/ui"
)

func Run() {

	if closed {
		log.Error("call to engine.Run ignored: ticker is closed")
		return
	}

	if running {
		log.Warn("call to engine.Run ignored: already running")
		return
	}

	running = true
	defer func() { running = false }()
	defer Close()

	// Calls ui.Tick synchronously to make sure everything is initialized before drawing
	ui.Tick(0)
	_ = pool.Frontend.Submit(tickLoop)
	_ = pool.Frontend.Submit(drawLoop)

	app.Main()
}
