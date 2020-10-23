package ui

import (
    "fmt"
    "gonum.org/v1/plot"
    "gonum.org/v1/plot/plotter"
    "gonum.org/v1/plot/vg"
    "log"
    "math"
    "math/rand"
    "os"
    "time"
    
    "github.com/ul-gaul/go-basestation/ui/plotting"
    "github.com/ul-gaul/go-basestation/ui/plotting/lines"
    "github.com/ul-gaul/go-basestation/ui/plotting/vgg"
    "github.com/ul-gaul/go-basestation/ui/widgets"
    "github.com/ul-gaul/go-basestation/utils"
    
    "gioui.org/app"
    "gioui.org/io/system"
    "gioui.org/layout"
    "gioui.org/op"
    "gioui.org/widget/material"
    
    "gioui.org/font/gofont"
)

func RunGioui() {
    defer os.Exit(0)
    window := app.NewWindow(app.Title("GAUL - Base Station"))
    defer window.Close()
    if err := loop(window); err != nil {
        log.Panicln(err)
    }
}

func loop(window *app.Window) error {
    th := material.NewTheme(gofont.Collection())
    
    
    drawer := vgg.New()
    chart, err := plot.New()
    utils.CheckErr(err)
    updateChart := plotting.NewChartUpdater(window, drawer, chart)
    chart.Add(lines.NewOriginLines(lines.LineWidth(vg.Points(1))))
    chart.Add(plotter.NewGrid())
    
    plotters := make([]*plotting.Plotter, 3)
    for i := range plotters {
        plotters[i], err = plotting.NewPlotter(
            plotting.WithStyleIdx(i),
            plotting.WithLegend(fmt.Sprintf("Line #%d", i+1)),
            plotting.WithData(squarePlot(i)...),
        )
        utils.CheckErr(err)
        plotters[i].ApplyTo(chart)
    }
    
    drawer.Redraw()
    
    /************************** **************************/
    var getTabContent = func(i int, title string) layout.Widget {
        return func(gtx layout.Context) layout.Dimensions {
            return layout.Center.Layout(gtx, drawer.ToWidget(chart))
        }
    }
    
    tabs := []widgets.TabChild{
        widgets.Tabbed("General", getTabContent(0, "GENERAL")),
        widgets.Tabbed("Motor", getTabContent(1, "MOTOR")),
    }
    
    
    /************************** **************************/
    var ops op.Ops
    tabBar := widgets.Tab(th)
    
    ticker := time.NewTicker(1 * time.Second)
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
                p.Append(randomPoints(1)...)
            }
            
            updateChart()
        }
    }
}

// randomPoints returns some random x, y points.
func randomPoints(n int) plotter.XYs {
    pts := make(plotter.XYs, n)
    last := rand.Float64()
    for i := range pts {
        last += math.Copysign(10*rand.Float64(), rand.Float64() - 0.5)
        pts[i].X = last
        pts[i].Y = pts[i].X + math.Copysign(10*rand.Float64(), rand.Float64() - 0.5)
    }
    return pts
}

// squarePlot *TEMP* génère des points
func squarePlot(i int) plotter.XYs {
    const b, d = 5.0, 2.0
    dist := b + (d * float64(i))
    xys := make(plotter.XYs, 6)
    xys[0] = plotter.XY{0, 0}
    
    var calcSign = func(v float64) float64 {
        // (-1) ** ( floor( (v%4) / 2 ) % 2 )
        return math.Mod(math.Pow(-1, math.Floor(math.Mod(v, 4)/2)), 2)
    }
    var calc = func(v float64) float64 {
        return calcSign(v) * dist
    }
    
    for j := 1; j < len(xys); j++ {
        xys[j].X = calc(float64(i) + float64(j) - 1)
        xys[j].Y = calc(float64(i) + float64(j))
    }
    
    last := &xys[len(xys)-1]
    if i%2 == 0 {
        last.Y += math.Copysign(1, last.Y)
    } else {
        last.X += math.Copysign(1, last.X)
    }
    
    return xys
}