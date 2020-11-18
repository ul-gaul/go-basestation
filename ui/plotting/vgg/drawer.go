package vgg

import (
    "gioui.org/layout"
    "gonum.org/v1/plot"
    "gonum.org/v1/plot/vg"
    "gonum.org/v1/plot/vg/draw"
    "image"
)

type Drawer struct {
    drawer draw.Canvas
    canvas *Canvas
    opts   []option
    redraw bool
    firstRun bool
    // TODO MinWidth, MinHeight, MaxWidth, MaxHeight
}

func New(opts ...option) *Drawer {
    canvas := NewCanvas(opts...)
    return &Drawer{
        drawer: draw.New(canvas),
        canvas: canvas,
        opts: opts,
        firstRun: true,
    }
}

func (d *Drawer) isSameSize(gtx layout.Context, min, max image.Point) bool {
    xmin, ymin, xmax, ymax :=
        d.toPixels(gtx, d.drawer.Min.X),
        d.toPixels(gtx, d.drawer.Min.Y),
        d.toPixels(gtx, d.drawer.Max.X),
        d.toPixels(gtx, d.drawer.Max.Y)
    
    return min.X == xmin && min.Y == ymin &&
        max.X == xmax && max.Y == ymax
}

func (d *Drawer) toPixels(gtx layout.Context, length vg.Length) int {
    return int(float64(gtx.Metric.PxPerSp) * length.Dots(d.canvas.DPI()))
}

func (d *Drawer) toLength(gtx layout.Context, px int) vg.Length {
    dpi := vg.Length(d.canvas.DPI())
    x := vg.Length(px)/vg.Length(gtx.Metric.PxPerSp)
    return (x/dpi) * vg.Inch
}

func (d *Drawer) Redraw() {
    d.redraw = true
}

func (d *Drawer) Layout(gtx layout.Context, chart *plot.Plot) layout.Dimensions {
    if d.firstRun || !d.isSameSize(gtx, gtx.Constraints.Min, gtx.Constraints.Max) {
        w := d.toLength(gtx, gtx.Constraints.Max.X) - d.toLength(gtx, gtx.Constraints.Min.X)
        h := d.toLength(gtx, gtx.Constraints.Max.Y) - d.toLength(gtx, gtx.Constraints.Min.Y)
        if w > 0 && h > 0 { // TODO Set to MinWidth & MinHeight
            d.canvas = NewCanvas(append([]option{UseWH(w, h)}, d.opts...)...)
            d.drawer = draw.New(d.canvas)
            d.firstRun = false
            d.redraw = true
        }
    }
    
    if d.redraw {
        chart.Draw(d.drawer)
        d.redraw = false
    }
    return d.canvas.Layout(gtx)
}


func (d *Drawer) ToWidget(chart *plot.Plot) layout.Widget {
    return func(gtx layout.Context) layout.Dimensions {
        return d.Layout(gtx, chart)
    }
}