package types

import (
	"gioui.org/io/key"
	"gioui.org/layout"
	"time"
)

type ITickable interface {
	Tick(delta time.Duration)
}

type IDrawable interface {
	Draw(gtx layout.Context) layout.Dimensions
}

type IComponent interface {
	ITickable
	IDrawable
}

type IView interface {
	IComponent
	Keypress(ev key.Event)
}