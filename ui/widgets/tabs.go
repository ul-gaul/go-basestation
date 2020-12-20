package widgets

import (
    "gioui.org/f32"
    "gioui.org/layout"
    "gioui.org/op/clip"
    "gioui.org/op/paint"
    "gioui.org/unit"
    "gioui.org/widget"
    "gioui.org/widget/material"
    "image"
)

type TabLayout struct {
    theme *material.Theme
    list  layout.List
    slider   TabSlider
    selected int
}

func Tab(th *material.Theme) *TabLayout {
    return &TabLayout{theme: th}
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

func (tb *TabLayout) Layout(gtx layout.Context, tabs ...TabChild) layout.Dimensions {
    if tb.selected >= len(tabs) {
        tb.selected = len(tabs) - 1
    }
    return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
        layout.Rigid(func(gtx layout.Context) layout.Dimensions {
            return tb.list.Layout(gtx, len(tabs), func(gtx layout.Context, tabIdx int) layout.Dimensions {
                t := &tabs[tabIdx]
                if t.btn.Clicked() {
                    if tb.selected < tabIdx {
                        tb.slider.PushLeft()
                    } else if tb.selected > tabIdx {
                        tb.slider.PushRight()
                    }
                    tb.selected = tabIdx
                }
                
                return material.Clickable(gtx, &t.btn, func(gtx layout.Context) layout.Dimensions {
                    return t.Button(gtx, tb.selected == tabIdx, tb.theme)
                })
            })
        }),
        // Contenu de chaque onglet
        layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
            return tb.slider.Layout(gtx, tabs[tb.selected].Content)
        }),
    )
}

func DefaultTabButton(gtx layout.Context, th *material.Theme, title string, isCurrent bool) layout.Dimensions {
    var tabWidth int
    return layout.Stack{Alignment: layout.S}.Layout(gtx,
        // Bouton et texte de chaque onglet
        layout.Stacked(func(gtx layout.Context) layout.Dimensions {
            dims := layout.UniformInset(unit.Sp(8)).Layout(gtx,
                material.Body1(th, title).Layout)
            tabWidth = dims.Size.X
            return dims
        }),
        // Barre sous l'onglet actif
        layout.Stacked(func(gtx layout.Context) layout.Dimensions {
            if !isCurrent {
                return layout.Dimensions{}
            }
            paint.ColorOp{Color: th.Fg}.Add(gtx.Ops)
            tabHeight := gtx.Px(unit.Dp(4))
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
