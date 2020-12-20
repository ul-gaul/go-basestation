package plotting

import (
    "gonum.org/v1/plot"
    "math"
    "reflect"
    
    "github.com/ul-gaul/go-basestation/utils"
)

const marginRatio = 0.05
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
    
    marginX, marginY := marginRatio * (xmax - xmin), marginRatio * (ymax - ymin)
    chart.X.Min, chart.X.Max = xmin - marginX, xmax + marginX
    chart.Y.Min, chart.Y.Max = ymin - marginY, ymax + marginY
}