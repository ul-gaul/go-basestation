package views

import (
    "gonum.org/v1/plot"
    
    "github.com/ul-gaul/go-basestation/ui/plotting"
    "github.com/ul-gaul/go-basestation/ui/plotting/vgg"
)

type PlotId uint8

const (
    PltAltitude PlotId = 0b0000_0001 << iota
    PltCoords
    PltTemperature
    PltPressure
)

type TabView interface {
    Plotters() map[PlotId]*plotting.Plotter
    Drawers() map[PlotId]*vgg.Drawer
    Charts() map[PlotId]*plot.Plot
}
