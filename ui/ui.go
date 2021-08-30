package ui

import (
	"fmt"
	"gioui.org/app"
	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/widget/material"
	"github.com/ul-gaul/go-basestation/cfg"
	"github.com/ul-gaul/go-basestation/ui/menu"
	"github.com/ul-gaul/go-basestation/ui/vars"
	"github.com/ul-gaul/go-basestation/ui/views"
	"strings"
	"time"
)

func Tick(delta time.Duration) {
	if !menu.Visible {
		views.CurrentView.Tick(delta)
	}
}

func Draw(gtx layout.Context, fps, tps int) layout.Dimensions {
	// Uncomment to print the size of the window
	// log.Printf("Size: %+v\n", gtx.Constraints.Max.Div(int(gtx.Metric.PxPerDp)))
	
	children := []layout.StackChild{
		layout.Stacked(views.CurrentView.Draw),
	}
	
	if cfg.Frontend.ShowFPS || cfg.Frontend.ShowTPS {
		children = append(children, layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			var str string
			
			if cfg.Frontend.ShowFPS {
				str += fmt.Sprintf("\n%d FPS", fps)
			}
			
			if cfg.Frontend.ShowTPS {
				str += fmt.Sprintf("\n%d TPS", tps)
			}
			
			str = strings.TrimSpace(str)
			
			txt := material.Body1(vars.Theme, str)
			txt.Alignment = text.End
			
			op.InvalidateOp{}.Add(gtx.Ops)
			return txt.Layout(gtx)
		}))
	}
	
	if menu.Visible {
		children = append(children, layout.Expanded(menu.Draw))
	}

	return layout.Stack{}.Layout(gtx, children...)
}

var isFullscreen bool
func Keypress(ev key.Event) {
	if ev.State == key.Press {
		skipViewHandler := false
		
		// Toggle Menu (ESC)
		if ev.Name == key.NameEscape && ev.Modifiers == 0 {
			skipViewHandler = true
			menu.Visible = !menu.Visible
		}
		
		// Toggle Fullscreen
		if ev.Name == "F" && ev.Modifiers == key.ModAlt {
			skipViewHandler = true
			isFullscreen = !isFullscreen
			
			if isFullscreen {
				vars.Window.Option(app.Windowed)
			} else {
				vars.Window.Option(app.Fullscreen)
			}
		}
		
		if !skipViewHandler {
			views.CurrentView.Keypress(ev)
		}
	}
}
