package ui

import (
    "fmt"
    "gioui.org/unit"
    log "github.com/sirupsen/logrus"
    "gonum.org/v1/plot/plotter"
    "os"
    "time"
    
    "github.com/ul-gaul/go-basestation/ui/plotting"
    "github.com/ul-gaul/go-basestation/ui/plotting/lines"
    ticker2 "github.com/ul-gaul/go-basestation/ui/plotting/ticker"
    "github.com/ul-gaul/go-basestation/ui/views"
    "github.com/ul-gaul/go-basestation/ui/widgets"
    "github.com/ul-gaul/go-basestation/utils"
    
    "gioui.org/app"
    "gioui.org/io/system"
    "gioui.org/layout"
    "gioui.org/op"
    "gioui.org/widget/material"
    
    "gioui.org/font/gofont"
)

var (
	window *app.Window
	th *material.Theme
)

func RunGioui() {
    defer os.Exit(0)
    window = app.NewWindow(app.Title("GAUL - Base Station"))
    th = material.NewTheme(gofont.Collection())
    
    defer window.Close()
    if err := loop(); err != nil {
        log.Panicln(err)
    }
}


func loop() error {
    
    generalTab, err := views.NewGeneralTab()
    if err != nil {
        return err
    }
    
    
    /************************** **************************/
    drawer, err := plotting.NewPlotDrawer()
    utils.CheckErr(err)
    drawer.Chart().Add(lines.NewOriginLines())
    drawer.Chart().Add(plotter.NewGrid())
    drawer.Chart().X.Tick.Marker = ticker2.NewTicker(10, ticker2.ContainData)
    
    plter, err := plotting.NewPlotter(
        plotting.WithStyleIdx(0),
        plotting.WithLegend(fmt.Sprintf("Line #%d", 1)),
        plotting.WithData(squarePlot(0)...),
    )
    utils.CheckErr(err)
    utils.CheckErr(drawer.AddPlotter(plter))
    
    
    tabs := []widgets.TabChild{
        widgets.Tabbed("General", generalTab.Layout),
        widgets.Tabbed("Motor", func(gtx layout.Context) layout.Dimensions {
            return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
                return layout.UniformInset(unit.Px(10)).Layout(gtx, drawer.Layout)
            })
        }),
    }
    
    /************************** **************************/
    ops := new(op.Ops)
    tabBar := widgets.Tab(th)
    
    ticker := time.NewTicker(150 * time.Millisecond)
    defer ticker.Stop()
    
    for {
        select {
        case e := <-window.Events():
            switch e := e.(type) {
            case system.DestroyEvent:
                return e.Err
            case system.FrameEvent:
                gtx := layout.NewContext(ops, e)
                tabBar.Layout(gtx, tabs...)
                e.Frame(gtx.Ops)
            }
        case <-ticker.C:
            addRandomPoints(plter, 1)
            log.Infof("Points: %d", plter.Data().Len())
        }
    }
}