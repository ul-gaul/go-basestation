package views

import (
    "gioui.org/layout"
    "gioui.org/unit"
    "gonum.org/v1/plot"
    "gonum.org/v1/plot/plotter"
    
    "github.com/ul-gaul/go-basestation/ui/plotting"
    "github.com/ul-gaul/go-basestation/ui/plotting/lines"
    "github.com/ul-gaul/go-basestation/ui/plotting/vgg"
)

var _ TabView = (*GeneralTab)(nil)

func NewGeneralTab() (g *GeneralTab, err error) {
    g = &GeneralTab{
        plotters: make(map[PlotId]*plotting.Plotter),
        drawers:  make(map[PlotId]*vgg.Drawer),
        charts:   make(map[PlotId]*plot.Plot),
    }
    
    if g.plotters[PltAltitude], g.charts[PltAltitude], err =
        createPlot("Altitude", "Time (ms)", "Height (m)"); err != nil {
        return
    }
    
    if g.plotters[PltCoords], g.charts[PltCoords], err =
        createPlot("Position", "Longitude", "Latitude"); err != nil {
        return
    }
    
    if g.plotters[PltTemperature], g.charts[PltTemperature], err =
        createPlot("Temperature", "Time (ms)", "Temperature (°C)"); err != nil {
        return
    }
    
    if g.plotters[PltPressure], g.charts[PltPressure], err =
        createPlot("Pressure", "Time (ms)", "??? (?)"); err != nil {
        return
    }
    
    for pltId := range g.plotters {
        g.drawers[pltId] = vgg.New()
    }
    
    return
}

type GeneralTab struct {
    plotters    map[PlotId]*plotting.Plotter
    drawers     map[PlotId]*vgg.Drawer
    charts      map[PlotId]*plot.Plot
}

func (g *GeneralTab) Plotters() map[PlotId]*plotting.Plotter { return g.plotters }
func (g *GeneralTab) Drawers() map[PlotId]*vgg.Drawer { return g.drawers }
func (g *GeneralTab) Charts() map[PlotId]*plot.Plot { return g.charts }

func createPlot(title, xAxis, yAxis string) (plt *plotting.Plotter, chart *plot.Plot, err error) {
    plt, err = plotting.NewPlotter()
    if err != nil {
        return
    }
    
    chart, err = plot.New()
    if err != nil {
        return
    }
    
    chart.Title.Text = title
    chart.X.Label.Text = xAxis
    chart.Y.Label.Text = yAxis
    chart.Add(lines.NewOriginLines())
    chart.Add(plotter.NewGrid())
    plt.ApplyTo(chart)
    return
}

func (g *GeneralTab) Layout(gtx layout.Context) layout.Dimensions {
    flexedChart := func(id PlotId) layout.FlexChild {
        return layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
            return layout.UniformInset(unit.Sp(10)).Layout(gtx,
                g.drawers[id].ToWidget(g.charts[id]))
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
