package menu

import (
	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/ul-gaul/go-basestation/ui/vars"
	"image/color"
)

var (
	btns []*Button
	list = layout.List{
		Axis:        layout.Vertical,
		ScrollToEnd: false,
		Alignment:   layout.Middle,
	}
	hideBtn widget.Clickable
	Visible bool
)

type Button struct {
	btn         widget.Clickable
	Label       string
	OnClick     func()
	VisibleCond func() bool
}

func AddButton(lbl string, callback func()) {
	AddButtonCond(lbl, func() bool { return true }, callback)
}

func AddButtonCond(lbl string, cond func() bool, callback func()) {
	btns = append(btns, &Button{
		Label:       lbl,
		OnClick:     callback,
		VisibleCond: cond,
	})
}

func Show() { Visible = true }
func Hide() { Visible = false }


func Draw(gtx layout.Context) layout.Dimensions {
	paint.Fill(gtx.Ops, color.NRGBA{R: 0x66, G: 0x66, B: 0x66, A: 0x72})
	if hideBtn.Clicked() {
		Hide()
	}
	hideBtn.Layout(gtx)
	
	return layout.Flex{
		Axis:      layout.Horizontal,
		Spacing:   layout.SpaceSides,
		Alignment: layout.Middle,
	}.Layout(gtx, layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{
			Axis:      layout.Vertical,
			Spacing:   layout.SpaceSides,
			Alignment: layout.Middle,
		}.Layout(gtx, layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return layout.Stack{Alignment: layout.Center}.Layout(gtx,
				layout.Expanded(func(gtx layout.Context) layout.Dimensions {
					w, h := float32(gtx.Constraints.Min.X), float32(gtx.Constraints.Min.Y)
					clip.UniformRRect(f32.Rect(-w/2, -h/2, w/2, h/2), 10).Add(gtx.Ops)
					paint.Fill(gtx.Ops, color.NRGBA{R: 0xAA, G: 0xAA, B: 0xAA, A: 0xFF})
					return layout.Dimensions{}
				}),
				layout.Stacked(func(gtx layout.Context) layout.Dimensions {
					var visibleBtns []*Button
					for _, btn := range btns {
						if btn.VisibleCond() {
							visibleBtns = append(visibleBtns, btn)
						}
					}
					
					return list.Layout(gtx, len(visibleBtns), func(gtx layout.Context, idx int) layout.Dimensions {
						mb := visibleBtns[idx]
						if mb.btn.Clicked() {
							mb.OnClick()
						}
						
						return layout.UniformInset(unit.Dp(16)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return material.Button(vars.Theme, &mb.btn, mb.Label).Layout(gtx)
						})
					})
				}))
		}))
	}))
}
