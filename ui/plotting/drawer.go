package plotting

import (
    "gioui.org/f32"
    "gioui.org/layout"
    "gioui.org/op"
    "gioui.org/op/paint"
    "gonum.org/v1/plot"
    "gonum.org/v1/plot/vg"
    "gonum.org/v1/plot/vg/draw"
    "gonum.org/v1/plot/vg/vgimg"
    "image"
    draw2 "image/draw"
    "sync"
    "time"
    
    "github.com/ul-gaul/go-basestation/pool"
)

type PlotDrawer struct {
    chart    *plot.Plot
    drawer   draw.Canvas
    canvas   *vgimg.Canvas
    plotters map[*Plotter]time.Time
    
    img draw2.Image
    
    mut         sync.RWMutex
    initialized bool
    chDraw      chan time.Time
}

// NewPlotDrawer creates a new PlotDrawer
func NewPlotDrawer(opts ...option) (*PlotDrawer, error) {
    var err error
    c := &PlotDrawer{
        plotters: make(map[*Plotter]time.Time),
        chDraw:   make(chan time.Time),
    }
    c.chart, err = plot.New()
    if err != nil {
        return nil, err
    }
    c.canvas = newCanvas(opts...)
    c.drawer = draw.New(c.canvas)
    c.img = c.cloneImg()
    return c, pool.Frontend.Submit(c.update)
}

// Chart returns the plot.Plot on which the plotters are drawn
func (c *PlotDrawer) Chart() *plot.Plot {
    return c.chart
}

// AddPlotter adds a Plotter to draw
func (c *PlotDrawer) AddPlotter(p *Plotter) error {
    var err error
    if _, ok := c.plotters[p]; !ok {
        c.plotters[p] = time.Now()
        c.chart.Add(p)
        if p.name != "" {
            c.chart.Legend.Add(p.name, p.line, p.points)
        }
        err = pool.Frontend.Submit(func() {
            for t := range p.chChange {
                if t.After(c.plotters[p]) {
                    c.plotters[p] = time.Now()
                    c.chDraw <- c.plotters[p]
                }
            }
        })
    }
    return err
}

// Layout renders the plot.Plot to the provided layout.Context
func (c *PlotDrawer) Layout(gtx layout.Context) layout.Dimensions {
    r32 := f32.Rect(
        float32(gtx.Constraints.Min.X),
        float32(gtx.Constraints.Min.Y),
        float32(gtx.Constraints.Max.X),
        float32(gtx.Constraints.Max.Y))
    
    c.mut.RLock()
    img := c.img
    c.mut.RUnlock()
    
    macro := op.Record(gtx.Ops)
    op.InvalidateOp{}.Add(gtx.Ops)
    paint.NewImageOp(img).Add(gtx.Ops)
    paint.PaintOp{Rect: r32}.Add(gtx.Ops)
    macro.Stop().Add(gtx.Ops)
    
    return layout.Dimensions{
        Size: image.Pt(int(r32.Max.X), int(r32.Max.Y)),
    }
}

func (c *PlotDrawer) update() {
    lastDraw := time.Now()
    doUpdate := func() {
        RecalcAxis(c.chart)
        c.chart.Draw(c.drawer)
        img := c.cloneImg()
        c.mut.Lock()
        c.img = img
        c.mut.Unlock()
    }
    
    doUpdate()
    for t := range c.chDraw {
        if t.After(lastDraw) {
            lastDraw = time.Now()
            doUpdate()
        }
    }
}

func (c *PlotDrawer) pt32(p vg.Point) f32.Point {
    _, h := c.canvas.Size()
    dpi := c.canvas.DPI()
    return f32.Point{
        X: float32(p.X.Dots(dpi)),
        Y: float32(h.Dots(dpi) - p.Y.Dots(dpi)),
    }
}

func (c *PlotDrawer) cloneImg() draw2.Image {
    src := c.canvas.Image()
    img := draw2.Image(image.NewRGBA(image.Rect(
        src.Bounds().Min.X, src.Bounds().Min.Y,
        src.Bounds().Max.X, src.Bounds().Max.Y)))
    draw2.Draw(img, img.Bounds(), image.Image(src), image.Point{}, draw2.Src)
    return img
}
