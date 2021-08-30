package engine

import (
	"gioui.org/io/key"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	log "github.com/sirupsen/logrus"
	"github.com/ul-gaul/go-basestation/ui"
	"github.com/ul-gaul/go-basestation/ui/vars"
	"time"
)

func tickLoop() {
	defer Close()

	lastTick := time.Now()
	for !closed {
		diff := time.Since(lastTick)
		if diff >= tickMinDeltaT {
			tps = int(time.Second / diff)
			lastTick = time.Now()
			ui.Tick(diff)
		}
	}
}

func drawLoop() {
	defer Close()

	var ops op.Ops
	lastFrame := time.Now()
	for !closed {
		select {
		case ev := <-vars.Window.Events():
			switch ev := ev.(type) {
			case system.DestroyEvent:
				if ev.Err != nil {
					log.Error(ev.Err)
				}
				return
			case key.Event:
				ui.Keypress(ev)
			case system.FrameEvent:
				gtx := layout.NewContext(&ops, ev)
				fps = int(time.Second / time.Since(lastFrame))
				lastFrame = gtx.Now

				ui.Draw(gtx, fps, TPS())
				ev.Frame(gtx.Ops)
			}
		}
	}
}
