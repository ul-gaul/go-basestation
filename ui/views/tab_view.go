package views

import (
    "github.com/ul-gaul/go-basestation/ui/plotting"
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
    Drawers() map[PlotId]*plotting.PlotDrawer
}
