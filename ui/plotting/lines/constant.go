package lines

import (
    "gonum.org/v1/plot"
    "gonum.org/v1/plot/vg"
    "gonum.org/v1/plot/vg/draw"
    "image/color"
)

var _ plot.Plotter = (*ConstantLine)(nil)

// DefaultConstantLineStyle is the default style for origin lines.
var DefaultConstantLineStyle = draw.LineStyle{
    Color: color.RGBA{0x33, 0x33, 0x33, 0x10},
    Width: vg.Points(0.25),
}

type Orientation uint8

const (
    Horizontal Orientation = 0
    Vertical   Orientation = 1
)

// ConstantLine implements the plot.Plotter interface, drawing
// a line representing a constant.
type ConstantLine struct {
    draw.LineStyle
    Value float64
    Orientation
}

// Plot implements the plot.Plotter interface.
func (cl *ConstantLine) Plot(cnv draw.Canvas, plt *plot.Plot) {
    if cl.Color == nil {
        return
    }
    
    trX, trY := plt.Transforms(&cnv)
    
    var (
        xmin = cnv.Min.X
        xmax = cnv.Max.X
        ymin = cnv.Min.Y
        ymax = cnv.Max.Y
    )
    
    switch cl.Orientation {
    case Horizontal:
        x := trX(cl.Value)
        if xmin < x && x < xmax {
            cnv.StrokeLine2(cl.LineStyle, x, ymin, x, ymax)
        }
    case Vertical:
        y := trY(cl.Value)
        if ymin < y && y < ymax {
            cnv.StrokeLine2(cl.LineStyle, xmin, y, xmax, y)
        }
    }
}

// NewConstantLine returns a new constant line.
func NewConstantLine(value float64, o Orientation, opts ...LineOption) *ConstantLine {
    lineStyle := DefaultConstantLineStyle
    for _, opt := range opts {
        opt(lineStyle)
    }
    
    return &ConstantLine{
        LineStyle:   lineStyle,
        Value:       value,
        Orientation: o,
    }
}

func NewOriginLine(o Orientation, opts ...LineOption) *ConstantLine {
    return NewConstantLine(0, o, opts...)
}

func NewOriginLines(opts ...LineOption) (xAxis, yAxis *ConstantLine) {
    return NewOriginLine(Horizontal, opts...), NewOriginLine(Vertical, opts...)
}