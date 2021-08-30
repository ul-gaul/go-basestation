package components

import (
	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"image"
	"math"
)

type TabLayout struct {
	Theme    *material.Theme
	Tabs     []TabChild
	list     layout.List
	selected int
}

func NewTabLayout(th *material.Theme, tabs ...TabChild) *TabLayout {
	return &TabLayout{
		Theme: th,
		Tabs:  tabs,
	}
}

type TabChild struct {
	btn     widget.Clickable
	Button  TabBarButton
	Content layout.Widget
}
type TabBarButton func(gtx layout.Context, isCurrent bool, theme *material.Theme) layout.Dimensions

func Tabbed(title string, content layout.Widget) TabChild {
	return TabChild{
		Content: content,
		Button: func(gtx layout.Context, isCurrent bool, theme *material.Theme) layout.Dimensions {
			return DefaultTabButton(gtx, theme, title, isCurrent)
		},
	}
}

func (tl *TabLayout) CurrentTab() int {
	return tl.selected
}

func (tl *TabLayout) SetCurrentTab(i int) {
	tl.selected = int(math.Max(0, math.Min(float64(len(tl.Tabs)-1), float64(i))))
}

func (tl *TabLayout) TabBar(gtx layout.Context) layout.Dimensions {
	return tl.list.Layout(gtx, len(tl.Tabs), func(gtx layout.Context, tabIdx int) layout.Dimensions {
		t := &tl.Tabs[tabIdx]
		if t.btn.Clicked() {
			tl.selected = tabIdx
		}

		return material.Clickable(gtx, &t.btn, func(gtx layout.Context) layout.Dimensions {
			return t.Button(gtx, tl.selected == tabIdx, tl.Theme)
		})
	})
}

func (tl *TabLayout) Content(gtx layout.Context) layout.Dimensions {
	return tl.Tabs[tl.selected].Content(gtx)
}

func (tl *TabLayout) Layout(gtx layout.Context) layout.Dimensions {
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(tl.TabBar),
		layout.Flexed(1, tl.Content),
	)
}

func DefaultTabButton(gtx layout.Context, th *material.Theme, title string, isCurrent bool) layout.Dimensions {
	var tabWidth int
	return layout.Stack{Alignment: layout.S}.Layout(gtx,
		// Bouton et texte de chaque onglet
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			dims := layout.Inset{
				Top:    unit.Sp(4),
				Right:  unit.Sp(8),
				Bottom: unit.Sp(4),
				Left:   unit.Sp(8),
			}.Layout(gtx, material.Body1(th, title).Layout)
			tabWidth = dims.Size.X
			return dims
		}),
		// Barre sous l'onglet actif
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			if !isCurrent {
				return layout.Dimensions{}
			}
			paint.ColorOp{Color: th.Fg}.Add(gtx.Ops)
			tabHeight := gtx.Px(unit.Dp(3))
			clip.UniformRRect(f32.Rectangle{
				Max: f32.Point{
					X: float32(tabWidth),
					Y: float32(tabHeight),
				},
			}, 5).Add(gtx.Ops)
			paint.PaintOp{}.Add(gtx.Ops)
			return layout.Dimensions{
				Size: image.Point{X: tabWidth, Y: tabHeight},
			}
		}),
	)
}
