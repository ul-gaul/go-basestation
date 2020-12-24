package plotting

import (
    "gioui.org/layout"
    "gioui.org/op"
    "gioui.org/op/clip"
    "gioui.org/op/paint"
    "github.com/disintegration/imaging"
    "gonum.org/v1/plot"
    "gonum.org/v1/plot/vg"
    "gonum.org/v1/plot/vg/draw"
    "gonum.org/v1/plot/vg/vgimg"
    "image"
    draw2 "image/draw"
    "math"
    "sync"
    "time"
    
    cfg "github.com/ul-gaul/go-basestation/config"
    "github.com/ul-gaul/go-basestation/constants"
    "github.com/ul-gaul/go-basestation/pool"
    "github.com/ul-gaul/go-basestation/ui/plotting/ticker"
)

const (
    DefaultWidth  vg.Length = 600
    DefaultHeight vg.Length = 400
    
    DefaultPadding = 0.03
)

type PlotDrawer struct {
    chart    *plot.Plot
    drawer   draw.Canvas
    canvas   *vgimg.Canvas
    plotters map[*Plotter]time.Time
    
    img draw2.Image
    
    w, h float64
    dpi  float64
    
    mut         sync.RWMutex
    initialized bool
    chDraw      chan struct{}
    
    paddingX, paddingY float64
}

// NewPlotDrawer creates a new PlotDrawer
func NewPlotDrawer() (*PlotDrawer, error) {
    var err error
    d := &PlotDrawer{
        plotters: make(map[*Plotter]time.Time),
        chDraw:   make(chan struct{}, 1),
        paddingX: DefaultPadding,
        paddingY: DefaultPadding,
    }
    d.chart, err = plot.New()
    if err != nil {
        return nil, err
    }
    
    d.chart.X.Tick.Marker = ticker.NewTicker(0, 0)
    d.chart.Y.Tick.Marker = ticker.NewTicker(0, 0)
    d.chart.X.Padding = 0
    d.chart.Y.Padding = 0
    
    d.canvas = vgimg.New(DefaultWidth, DefaultHeight)
    d.dpi = d.canvas.DPI()
    d.w, d.h = DefaultWidth.Dots(d.dpi), DefaultHeight.Dots(d.dpi)
    
    d.drawer = draw.New(d.canvas)
    d.img = imaging.Clone(d.canvas.Image())
    
    return d, pool.Frontend.Submit(d.update)
}

// Chart returns the plot.Plot on which the plotters are drawn
func (d *PlotDrawer) Chart() *plot.Plot {
    return d.chart
}

// AddPlotter adds a Plotter to draw
func (d *PlotDrawer) AddPlotter(p *Plotter) error {
    var err error
    if _, ok := d.plotters[p]; !ok {
        d.plotters[p] = time.Now()
        d.chart.Add(p)
        if p.name != "" {
            d.chart.Legend.Add(p.name, p.line, p.points)
        }
        err = pool.Frontend.Submit(func() {
            for t := range p.chChange {
                if t.After(d.plotters[p]) {
                    d.plotters[p] = time.Now()
                    d.redraw()
                }
            }
        })
    }
    return err
}

// Layout renders the plot.Plot to the provided layout.Context
func (d *PlotDrawer) Layout(gtx layout.Context) layout.Dimensions {
    rect := image.Rect(
        0, 0,
        gtx.Constraints.Max.X,
        gtx.Constraints.Max.Y)
    
    // 1 dp = 160 dpi
    dpi := float64(gtx.Metric.PxPerDp) * 160 * cfg.Frontend.PlotScale
    w, h := float64(rect.Dx()), float64(rect.Dy())
    if (dpi > 0 && w > 0 && h > 0) && (dpi != d.dpi || w != d.w || h != d.h) {
        d.dpi, d.w, d.h = dpi, w, h
        d.redraw()
    }
    
    macro := op.Record(gtx.Ops)
    op.InvalidateOp{}.Add(gtx.Ops)
    paint.NewImageOp(d.img).Add(gtx.Ops)
    clip.Rect(rect).Add(gtx.Ops)
    paint.PaintOp{}.Add(gtx.Ops)
    macro.Stop().Add(gtx.Ops)
    
    return layout.Dimensions{
        Size: image.Pt(rect.Max.X, rect.Max.Y),
    }
}

func (d *PlotDrawer) redraw() {
    select {
    case d.chDraw <- struct{}{}:
    default:
    }
}

/************************** UPDATE **************************/

func (d *PlotDrawer) update() {
    d.canvas.Push()
    for range d.chDraw {
        cdpi := d.canvas.DPI()
        cw, ch := d.canvas.Size()
    
        d.recalcAxis(d.chart)
        if d.dpi != cdpi || d.w != cw.Dots(d.dpi) || d.h != ch.Dots(d.dpi) {
            d.canvas = vgimg.NewWith(
                vgimg.UseWH(vg.Length(d.w/d.dpi)*vg.Inch, vg.Length(d.h/d.dpi)*vg.Inch),
                vgimg.UseDPI(int(d.dpi)))
            d.drawer = draw.New(d.canvas)
        }
        d.chart.Draw(d.drawer)
        d.chart.DataCanvas(d.drawer)
        
        img := imaging.Clone(d.canvas.Image())
        d.img = img
    }
}


/************************** AXIS **************************/

func (d *PlotDrawer) recalcAxis(chart *plot.Plot) {
    var xmin, xmax, ymin, ymax float64
    for p := range d.plotters {
        pxmin, pxmax, pymin, pymax := p.DataRange()
        xmin = math.Min(xmin, pxmin)
        xmax = math.Max(xmax, pxmax)
        ymin = math.Min(ymin, pymin)
        ymax = math.Max(ymax, pymax)
    }
    
    xpad, ypad := d.paddingX*(xmax-xmin), d.paddingY*(ymax-ymin)
    chart.X.Min, chart.X.Max = xmin-xpad, xmax+xpad
    chart.Y.Min, chart.Y.Max = ymin-ypad, ymax+ypad
}

// PaddingX returns the padding ratio for the X axis.
func (d *PlotDrawer) PaddingX() float64 {
    return d.paddingX
}

// PaddingY returns the padding ratio for the Y axis.
func (d *PlotDrawer) PaddingY() float64 {
    return d.paddingY
}

// SetPaddingX sets the padding ratio for the X axis.
// Must be between -1 and 1;
func (d *PlotDrawer) SetPaddingX(padding float64) error {
    return d.setPadding(padding, d.paddingY)
}

// SetPaddingX sets the data padding for the Y axis.
// Must be between -1 and 1;
func (d *PlotDrawer) SetPaddingY(padding float64) error {
    return d.setPadding(d.paddingX, padding)
}

func (d *PlotDrawer) setPadding(x, y float64) error {
    if x < -1 || x > 1 || y < -1 || y > 1 {
        return constants.ErrPaddingOutOfRange
    }
    d.paddingX, d.paddingY = x, y
    d.redraw()
    return nil
}
