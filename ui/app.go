package ui

import (
    "fmt"
    "gioui.org/unit"
    log "github.com/sirupsen/logrus"
    "gonum.org/v1/plot"
    "gonum.org/v1/plot/plotter"
    "os"
    "reflect"
    "time"
    
    "github.com/ul-gaul/go-basestation/pool"
    "github.com/ul-gaul/go-basestation/ui/plotting"
    "github.com/ul-gaul/go-basestation/ui/plotting/lines"
    "github.com/ul-gaul/go-basestation/ui/plotting/vgg"
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

var window *app.Window

func RunGioui() {
    defer os.Exit(0)
    window = app.NewWindow(app.Title("GAUL - Base Station"))
    
    defer window.Close()
    if err := loop(); err != nil {
        log.Panicln(err)
    }
}

func startListening(tabs ...views.TabView) {
    var updaters []plotting.ChartUpdater
    var cases []reflect.SelectCase
    
    for _, tab := range tabs {
        for pltId, plt := range tab.Plotters() {
            updaters = append(updaters, plotting.NewChartUpdater(window,
                tab.Drawers()[pltId], tab.Charts()[pltId]))
            
            cases = append(cases, reflect.SelectCase{
                Dir:  reflect.SelectRecv,
                Chan: reflect.ValueOf(plt.ChangeChannel()),
            })
        }
    }
    
    // t := time.NewTicker(time.Second / 10)
    // defer t.Stop()
    // for {
    //     <- t.C
    //     for _, update := range updaters {
    //         update()
    //     }
    //     window.Invalidate()
    // }
    
    for {
        if idx, _, ok := reflect.Select(cases); ok {
            updaters[idx]()
        } else {
            log.Warnln("Listener stopped: channel closed!")
            return
        }
    }
}

func refreshLoop() {
    ticker := time.NewTicker(time.Second / 10)
    for {
        <- ticker.C
        window.Invalidate()
    }
}


func loop() error {
    th := material.NewTheme(gofont.Collection())
    
    generalTab, err := views.NewGeneralTab()
    if err != nil {
        return err
    }
    
    if err = pool.Frontend.Submit(func() { startListening(generalTab) }); err != nil {
        return err
    }
    if err = pool.Frontend.Submit(refreshLoop); err != nil {
        return err
    }
    
    /************************** **************************/
    drawer := vgg.New()
    chart, err := plot.New()
    utils.CheckErr(err)
    updateChart := plotting.NewChartUpdater(window, drawer, chart)
    chart.Add(lines.NewOriginLines())
    chart.Add(plotter.NewGrid())
    
    plotters := make([]*plotting.Plotter, 1)
    for i := range plotters {
        plotters[i], err = plotting.NewPlotter(
            plotting.WithStyleIdx(i),
            plotting.WithLegend(fmt.Sprintf("Line #%d", i+1)),
            plotting.WithData(squarePlot(i)...),
        )
        utils.CheckErr(err)
        plotters[i].ApplyTo(chart)
    }
    
    tabs := []widgets.TabChild{
        widgets.Tabbed("General", generalTab.Layout),
        widgets.Tabbed("Motor", func(gtx layout.Context) layout.Dimensions {
            return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
                return layout.UniformInset(unit.Sp(10)).Layout(gtx, drawer.ToWidget(chart))
            })
        }),
    }
    
    /************************** **************************/
    var ops op.Ops
    tabBar := widgets.Tab(th)
    
    ticker := time.NewTicker(time.Second / 10)
    defer ticker.Stop()
    
    for {
        select {
        case e := <-window.Events():
            switch e := e.(type) {
            case system.DestroyEvent:
                return e.Err
            case system.FrameEvent:
                gtx := layout.NewContext(&ops, e)
                tabBar.Layout(gtx, tabs...)
                
                e.Frame(gtx.Ops)
            }
        
        case <-ticker.C:
            for _, p := range plotters {
                addRandomPoints(p, 1)
            }
            
            log.Infof("Points: %d", plotters[0].Data().Len())
            
            updateChart()
        }
    }
}