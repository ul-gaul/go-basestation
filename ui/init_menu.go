package ui

import (
	"github.com/ul-gaul/go-basestation/controller"
	"github.com/ul-gaul/go-basestation/ui/menu"
	"github.com/ul-gaul/go-basestation/ui/vars"
	"github.com/ul-gaul/go-basestation/ui/views"
	"os"
)

func init() {
	modeCond := func(mode controller.Mode, conds ...func() bool) func() bool {
		return func() bool {
			visible := controller.CurrentMode() == mode
			for i := 0; i < len(conds) && visible; i++ {
				visible = conds[i]()
			}
			return visible
		}
	}

	not := func(cond func() bool) func() bool {
		return func() bool {
			return !cond()
		}
	}

	menu.AddButton("Resume", menu.Hide)

	// Remote/Serial button
	menu.AddButtonCond("Disconnect & Change Mode", modeCond(controller.MODE_SERIAL), func() {
		controller.CloseConnection()

		menu.Hide()
		views.CurrentView = views.Launcher
		vars.DataRoller.Clear()
	})

	// Replay buttons
	menu.AddButtonCond("Resume",
		modeCond(controller.MODE_REPLAY, controller.IsReplayPaused),
		controller.ResumeReplay)

	menu.AddButtonCond("Pause",
		modeCond(controller.MODE_REPLAY, not(controller.IsReplayPaused)),
		controller.PauseReplay)

	menu.AddButtonCond("Stop & Change Mode", modeCond(controller.MODE_REPLAY), func() {
		controller.StopReplay()
		menu.Hide()
		views.CurrentView = views.Launcher
		vars.DataRoller.Clear()
	})

	// Generate button
	menu.AddButtonCond("Stop & Change Mode", modeCond(controller.MODE_GENERATE), func() {
		controller.StopGenerator()
		menu.Hide()
		views.CurrentView = views.Launcher
		vars.DataRoller.Clear()
	})

	menu.AddButton("Exit", func() { os.Exit(0) })
}
