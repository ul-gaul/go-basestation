package views

import (
    "gioui.org/layout"
    "gioui.org/unit"
    "gonum.org/v1/plot/plotter"
    
    "github.com/ul-gaul/go-basestation/ui/plotting"
    "github.com/ul-gaul/go-basestation/ui/plotting/lines"
)

var _ TabView = (*GeneralTab)(nil)

func NewGeneralTab() (g *GeneralTab, err error) {
    g = &GeneralTab{
        plotters: make(map[PlotId]*plotting.Plotter),
        drawers:  make(map[PlotId]*plotting.PlotDrawer),
    }
    
    if g.plotters[PltAltitude], g.drawers[PltAltitude], err =
        createPlot("Altitude", "Time (ms)", "Height (m)"); err != nil {
        return
    }
    
    if g.plotters[PltCoords], g.drawers[PltCoords], err =
        createPlot("Position", "Longitude", "Latitude"); err != nil {
        return
    }
    
    if g.plotters[PltTemperature], g.drawers[PltTemperature], err =
        createPlot("Temperature", "Time (ms)", "Temperature (°C)"); err != nil {
        return
    }
    
    if g.plotters[PltPressure], g.drawers[PltPressure], err =
        createPlot("Pressure", "Time (ms)", "??? (?)"); err != nil {
        return
    }
    
    return
}

type GeneralTab struct {
    plotters map[PlotId]*plotting.Plotter
    drawers  map[PlotId]*plotting.PlotDrawer
}

func (g *GeneralTab) Plotters() map[PlotId]*plotting.Plotter   { return g.plotters }
func (g *GeneralTab) Drawers() map[PlotId]*plotting.PlotDrawer { return g.drawers }

func createPlot(title, xAxis, yAxis string) (plt *plotting.Plotter, drawer *plotting.PlotDrawer, err error) {
    plt, err = plotting.NewPlotter()
    if err != nil {
        return
    }
    
    drawer, err = plotting.NewPlotDrawer()
    if err != nil {
        return
    }
    
    chart := drawer.Chart()
    chart.Title.Text = title
    chart.X.Label.Text = xAxis
    chart.Y.Label.Text = yAxis
    chart.Add(lines.NewOriginLines())
    chart.Add(plotter.NewGrid())
    err = drawer.AddPlotter(plt)
    return
}

func (g *GeneralTab) Layout(gtx layout.Context) layout.Dimensions {
    flexedChart := func(id PlotId) layout.FlexChild {
        return layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
            return layout.UniformInset(unit.Px(10)).Layout(gtx, g.drawers[id].Layout)
        })
    }
    
    return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
        // Column 1
        layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
            return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
                // Col 1, Row 1
                flexedChart(PltAltitude),
                // Col 1, Row 2
                flexedChart(PltCoords),
            )
        }),
        
        // Column 2
        layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
            return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
                // Col 2, Row 1
                flexedChart(PltTemperature),
                // Col 2, Row 2
                flexedChart(PltPressure),
            )
        }),
    )
}
