package views

import (
	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	log "github.com/sirupsen/logrus"
	"github.com/sqweek/dialog"
	"github.com/ul-gaul/go-basestation/controller"
	"github.com/ul-gaul/go-basestation/ui/types"
	"github.com/ul-gaul/go-basestation/ui/vars"
	"image"
	"image/color"
	"os"
	"path/filepath"
	"time"
)

// TODO documentation

var _ types.IView = (*launcherView)(nil)

type launcherView struct {
	lastRefresh time.Time
	ports       []string
	portEnum    widget.Enum
	customPort  widget.Editor
	outputBtn   widget.Clickable
	connectBtn  widget.Clickable

	replayFileBtn widget.Clickable
	replayFile    *os.File
	replayBtn     widget.Clickable

	generateBtn widget.Clickable
}

func (l *launcherView) Keypress(_ key.Event) {}

func (l *launcherView) Tick(delta time.Duration) {
	// Refresh ports list every second
	if time.Since(l.lastRefresh) >= time.Second {
		l.ports = controller.ListAvailablePorts()
		l.lastRefresh = time.Now()

		// Check if selected port is still in the list
		found := false
		for i := 0; i < len(l.ports) && !found; i++ {
			found = l.ports[i] == l.portEnum.Value
		}

		if !found {
			l.portEnum.Value = ""
		}
	}
}

func (l *launcherView) Draw(gtx layout.Context) layout.Dimensions {
	if l.connectBtn.Clicked() {
		l.onConnect()
	}

	if l.replayBtn.Clicked() {
		l.onReplay()
	}

	if l.generateBtn.Clicked() {
		l.onGenerate()
	}

	return layout.Inset{Left: unit.Dp(20), Right: unit.Dp(20)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{
			Axis:      layout.Vertical,
			Spacing:   layout.SpaceEvenly,
			Alignment: layout.Middle,
		}.Layout(gtx,
			// Remote Connection
			layout.Flexed(1.2, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{
					Axis:      layout.Vertical,
					Spacing:   layout.SpaceSides,
					Alignment: layout.Middle,
				}.Layout(gtx,
					layout.Rigid(l.title("Remote / Serial")),
					layout.Rigid(l.spacer),
					layout.Rigid(l.portSelection),
					layout.Rigid(l.spacer),
					layout.Rigid(l.selectSerialOutput),
					layout.Rigid(l.spacer),
					layout.Rigid(material.Button(vars.Theme, &l.connectBtn, "Start Listenning").Layout),
				)
			}),

			layout.Rigid(l.divider),

			// Replay
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{
					Axis:      layout.Vertical,
					Spacing:   layout.SpaceSides,
					Alignment: layout.Middle,
				}.Layout(gtx,
					layout.Rigid(l.title("Replay")),
					layout.Rigid(l.spacer),
					layout.Rigid(l.selectReplayFile),
					layout.Rigid(l.spacer),
					layout.Rigid(material.Button(vars.Theme, &l.replayBtn, "Replay").Layout),
				)
			}),

			layout.Rigid(l.divider),

			// Generate
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{
					Axis:      layout.Vertical,
					Spacing:   layout.SpaceSides,
					Alignment: layout.Middle,
				}.Layout(gtx,
					layout.Rigid(l.title("Generate")),
					layout.Rigid(l.spacer),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						desc := material.Body2(vars.Theme, "Generates fake/random data")
						desc.Alignment = text.Middle
						return desc.Layout(gtx)
					}),
					layout.Rigid(l.spacer),
					layout.Rigid(material.Button(vars.Theme, &l.generateBtn, "Generate").Layout),
				)
			}),
		)
	})
}

func (l *launcherView) onConnect() {
	port := l.portEnum.Value
	if port == "" {
		port = l.customPort.Text()
	}

	err := controller.OpenConnection(port)
	if err == nil {
		if l.replayFile != nil {
			_ = l.replayFile.Close()
			l.replayFile = nil
		}
		CurrentView = Main
	} else {
		log.Error(err)
		dialog.Message("Could not connect to specified port:\n%v", err).Title("Error!").Error()
	}
}

func (l *launcherView) onReplay() {
	if l.replayFile != nil {
		if err := controller.Replay(l.replayFile); err == nil {
			CurrentView = Main
		} else {
			log.Error(err)
			dialog.Message("Could not start simulation:\n%v", err).Title("Error!").Error()
		}
	} else {
		dialog.Message("No file selected!").Title("Error!").Error()
	}
}

func (l *launcherView) onGenerate() {
	if err := controller.Generate(); err == nil {
		if l.replayFile != nil {
			_ = l.replayFile.Close()
			l.replayFile = nil
		}
		CurrentView = Main
	} else {
		dialog.Message("Could not start generator:\n%v", err).Title("Error!").Error()
	}
}

func (l *launcherView) spacer(gtx layout.Context) layout.Dimensions {
	return layout.Spacer{Height: unit.Dp(10)}.Layout(gtx)
}

func (l *launcherView) divider(gtx layout.Context) layout.Dimensions {
	w := float32(gtx.Constraints.Max.X)
	size := image.Pt(int(w*0.6), 1)
	// op.Offset(f32.Pt(w*0.1, 0)).Add(gtx.Ops)
	clip.Rect{Max: size}.Add(gtx.Ops)
	paint.Fill(gtx.Ops, color.NRGBA{A: 0xFF})
	return layout.Dimensions{Size: size}
}

func (l *launcherView) title(title string) layout.Widget {
	t := material.H5(vars.Theme, title)
	t.Alignment = text.Middle
	return t.Layout
}

func (l *launcherView) label(lbl string) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		return layout.Inset{Right: unit.Dp(10)}.Layout(gtx,
			material.Body1(vars.Theme, lbl).Layout)
	}
}

func (l *launcherView) portSelection(gtx layout.Context) layout.Dimensions {
	return layout.Flex{
		Axis:      layout.Horizontal,
		Spacing:   layout.SpaceBetween,
		Alignment: layout.Start,
	}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Top: unit.Dp(8), Right: unit.Dp(10)}.Layout(gtx,
				material.Body1(vars.Theme, "Serial Port :").Layout)
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			children := make([]layout.FlexChild, len(l.ports)+1)

			children[0] = layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				var dims layout.Dimensions
				if l.portEnum.Value != "" {
					radio := material.RadioButton(vars.Theme, &l.portEnum, "", "Custom Port")
					radio.Font.Style = text.Italic
					dims = radio.Layout(gtx)
				} else {
					dims = layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
						layout.Rigid(material.RadioButton(vars.Theme, &l.portEnum, "", "").Layout),
						layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
							return widget.Border{
								Color:        vars.Theme.ContrastBg,
								CornerRadius: unit.Dp(3),
								Width:        unit.Dp(1),
							}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
								return layout.UniformInset(unit.Dp(8)).Layout(gtx,
									material.Editor(vars.Theme, &l.customPort, "Ex: /dev/tty2").Layout)
							})
						}),
					)
				}
				return dims
			})

			for i, port := range l.ports {
				children[i+1] = layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					rb := material.RadioButton(vars.Theme, &l.portEnum, port, port)
					rb.Font.Variant = "Mono"
					return rb.Layout(gtx)
				})
			}

			return layout.Flex{
				Axis:      layout.Vertical,
				Spacing:   layout.SpaceEnd,
				Alignment: layout.Start,
			}.Layout(gtx, children...)
		}),
	)
}

func (l *launcherView) selectSerialOutput(gtx layout.Context) layout.Dimensions {
	if l.outputBtn.Clicked() {
		dir := filepath.Dir(controller.GetSerialOutputFile().Name())
		path, err := dialog.File().Filter("Comma Separated Value (CSV)", "csv").SetStartDir(dir).Save()
		if err == nil {
			err = controller.SetSerialOutputFile(path, false)
			if err == os.ErrExist {
				err = nil
				if dialog.Message("Already exists:\n%s\n\nOverwrite existing file?", path).Title("Overwrite?").YesNo() {
					err = controller.SetSerialOutputFile(path, true)
				}
			}
		}

		if err != nil {
			log.Error(err)
			dialog.Message("%v", err).Title("Error!").Error()
		}
	}

	return layout.Flex{
		Axis:      layout.Horizontal,
		Alignment: layout.Middle,
	}.Layout(gtx,
		layout.Rigid(l.label("Output : ")),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			txt := material.Body2(vars.Theme, controller.GetSerialOutputFile().Name())
			txt.Font.Variant = "Mono"
			return txt.Layout(gtx)
		}),
		layout.Rigid(material.Button(vars.Theme, &l.outputBtn, "Change").Layout),
	)
}

func (l *launcherView) selectReplayFile(gtx layout.Context) layout.Dimensions {
	if l.replayFileBtn.Clicked() {
		path, err := dialog.File().Filter("Comma Separated Values (CSV)", "csv").Load()
		if err != nil {
			log.Error(err)
			dialog.Message("Error opening file:\n%v", err).Title("Error!").Error()
		}
		l.replayFile, err = os.OpenFile(path, os.O_RDONLY, os.ModePerm)
		if err != nil {
			l.replayFile = nil
			log.Error(err)
			dialog.Message("Error opening file:\n%v", err).Title("Error!").Error()
		}
	}

	return layout.Flex{
		Axis:      layout.Horizontal,
		Alignment: layout.Middle,
	}.Layout(gtx,
		layout.Rigid(l.label("Data file : ")),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			lbl := material.Body2(vars.Theme, "<no file selected>")
			if l.replayFile == nil {
				lbl.Font.Style = text.Italic
			} else {
				lbl.Text = l.replayFile.Name()
				lbl.Font.Variant = "Mono"
			}
			return lbl.Layout(gtx)
		}),
		layout.Rigid(material.Button(vars.Theme, &l.replayFileBtn, "Choose").Layout),
	)
}
