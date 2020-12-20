package plotting

import (
    "gioui.org/layout"
    "gioui.org/op"
    "gioui.org/op/clip"
    "gioui.org/op/paint"
    "gonum.org/v1/plot"
    "gonum.org/v1/plot/vg"
    "gonum.org/v1/plot/vg/draw"
    "gonum.org/v1/plot/vg/vgimg"
    "image"
    draw2 "image/draw"
    "math"
    "sync"
    "time"
    
    "github.com/ul-gaul/go-basestation/constants"
    "github.com/ul-gaul/go-basestation/pool"
    "github.com/ul-gaul/go-basestation/ui/plotting/ticker"
)

const (
    DefaultWidth vg.Length = 600
    DefaultHeight vg.Length = 400
    
    DefaultPadding = 0.05
)

type PlotDrawer struct {
    chart    *plot.Plot
    drawer   draw.Canvas
    canvas   *vgimg.Canvas
    plotters map[*Plotter]time.Time
    
    img draw2.Image
    
    w, h float64
    dpi float64
    
    mut         sync.RWMutex
    initialized bool
    chDraw      chan time.Time
    
    paddingX, paddingY float64
}

// NewPlotDrawer creates a new PlotDrawer
func NewPlotDrawer() (*PlotDrawer, error) {
    var err error
    d := &PlotDrawer{
        plotters: make(map[*Plotter]time.Time),
        chDraw:   make(chan time.Time),
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
    d.img = d.cloneImg()
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
                    d.chDraw <- d.plotters[p]
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
    dpi := float64(gtx.Metric.PxPerDp) * 160
    w, h := float64(rect.Dx()), float64(rect.Dy())
    if (dpi > 0 && w > 0 && h > 0) && (dpi != d.dpi || w != d.w || h != d.h) {
        d.mut.Lock()
        d.dpi, d.w, d.h = dpi, w, h
        d.mut.Unlock()
        d.chDraw <- time.Now()
    }
    
    d.mut.RLock()
    img := d.img
    d.mut.RUnlock()
    
    macro := op.Record(gtx.Ops)
    op.InvalidateOp{}.Add(gtx.Ops)
    paint.NewImageOp(img).Add(gtx.Ops)
    clip.Rect(rect).Add(gtx.Ops)
    paint.PaintOp{}.Add(gtx.Ops)
    macro.Stop().Add(gtx.Ops)
    
    return layout.Dimensions{
        Size: image.Pt(rect.Max.X, rect.Max.Y),
    }
}

/************************** UPDATE **************************/

func (d *PlotDrawer) update() {
    lastDraw := time.Now()
    doUpdate := func() {
        d.mut.RLock()
        dpi, w, h := d.dpi, d.w, d.h
        d.mut.RUnlock()
        cdpi := d.canvas.DPI()
        cw, ch := d.canvas.Size()
        
        d.recalcAxis(d.chart)
        if dpi != cdpi || w != cw.Dots(dpi) || h != ch.Dots(dpi) {
            d.canvas = vgimg.NewWith(
                vgimg.UseWH(vg.Length(w/dpi) * vg.Inch, vg.Length(h/dpi) * vg.Inch),
                vgimg.UseDPI(int(dpi)))
            d.drawer = draw.New(d.canvas)
        }
        d.chart.Draw(d.drawer)
        img := d.cloneImg()
        
        d.mut.Lock()
        d.img = img
        d.mut.Unlock()
    }
    
    doUpdate()
    for t := range d.chDraw {
        if t.After(lastDraw) {
            lastDraw = time.Now()
            doUpdate()
        }
    }
}

func (d *PlotDrawer) cloneImg() draw2.Image {
    src := d.canvas.Image()
    img := draw2.Image(image.NewRGBA(image.Rect(
        src.Bounds().Min.X, src.Bounds().Min.Y,
        src.Bounds().Max.X, src.Bounds().Max.Y)))
    draw2.Draw(img, img.Bounds(), image.Image(src), image.Point{}, draw2.Src)
    return img
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
    
    xpad, ypad := d.paddingX * (xmax - xmin), d.paddingY * (ymax - ymin)
    chart.X.Min, chart.X.Max = xmin - xpad, xmax + xpad
    chart.Y.Min, chart.Y.Max = ymin - ypad, ymax + ypad
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
    d.mut.Lock()
    d.paddingX, d.paddingY = x, y
    d.mut.Unlock()
    d.chDraw <- time.Now()
    return nil
}
