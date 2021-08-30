package views

import (
	"gioui.org/io/key"
	"gioui.org/text"
	"gioui.org/widget"
	"github.com/ul-gaul/go-basestation/controller"
	"github.com/ul-gaul/go-basestation/ui/components"
	"github.com/ul-gaul/go-basestation/ui/types"
	"github.com/ul-gaul/go-basestation/ui/vars"
	tabs2 "github.com/ul-gaul/go-basestation/ui/views/tabs"
)

var (
	Main     *mainView
	Launcher *launcherView
)

var CurrentView types.IView

func init() {
	Main = &mainView{
		generalTab: tabs2.NewGeneralTab(),
		motorTab:   tabs2.NewMotorTab(),
	}
	Main.tabBar = components.NewTabLayout(vars.Theme,
		components.Tabbed("General", Main.generalTab.Draw),
		components.Tabbed("Motor", Main.motorTab.Draw))

	Launcher = &launcherView{
		ports: controller.ListAvailablePorts(),
		customPort: widget.Editor{
			Alignment:  text.Start,
			SingleLine: true,
			Submit:     true,
			InputHint:  key.HintText,
		},
	}

	CurrentView = Launcher
}
