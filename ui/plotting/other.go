package plotting

import (
    "gioui.org/app"
    "gonum.org/v1/plot"
    "math"
    "reflect"
    
    "github.com/ul-gaul/go-basestation/ui/plotting/vgg"
    "github.com/ul-gaul/go-basestation/utils"
)

func RecalcAxis(chart *plot.Plot) {
    v := reflect.Indirect(reflect.ValueOf(chart)).FieldByName("plotters")
    plotters := utils.GetUnexportedField(v).([]plot.Plotter)
    
    var xmin, xmax, ymin, ymax float64
    for _, p := range plotters {
        if dr, ok := p.(plot.DataRanger); ok {
            pxmin, pxmax, pymin, pymax := dr.DataRange()
            xmin = math.Min(xmin, pxmin)
            xmax = math.Max(xmax, pxmax)
            ymin = math.Min(ymin, pymin)
            ymax = math.Max(ymax, pymax)
        }
    }
    
    chart.X.Min, chart.X.Max = xmin, xmax
    chart.Y.Min, chart.Y.Max = ymin, ymax
}

type ChartUpdater func()

func NewChartUpdater(win *app.Window, drawer *vgg.Drawer, chart *plot.Plot) ChartUpdater {
    return func() {
        RecalcAxis(chart)
        drawer.Redraw()
        win.Invalidate()
    }
}
