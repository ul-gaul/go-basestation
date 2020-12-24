package plotting

import (
    "gonum.org/v1/plot"
    "gonum.org/v1/plot/plotter"
    "gonum.org/v1/plot/vg/draw"
    "sync"
    "time"
    
    "github.com/ul-gaul/go-basestation/constants"
    "github.com/ul-gaul/go-basestation/utils"
)

// TODO Documentation

var _ plot.Plotter = (*Plotter)(nil)

type Plotter struct {
    Name                   string
    LineStyle              draw.LineStyle
    PointStyleFunc         func(int) draw.GlyphStyle
    DataLimit              int
    xys                    plotter.XYs
    padRatioX, padRatioY   float64
    chChange               chan time.Time
    mut                    sync.Mutex
}

func (p *Plotter) Data() plotter.XYs { return p.xys }

func (p *Plotter) Prepend(xys ...plotter.XY)       { p.PrependAll(xys) }
func (p *Plotter) PrependAll(xys plotter.XYs)      { p.InsertAll(0, xys) }
func (p *Plotter) Append(xys ...plotter.XY)        { p.AppendAll(xys) }
func (p *Plotter) AppendAll(xys plotter.XYs)       { p.InsertAll(len(p.xys), xys) }
func (p *Plotter) Insert(i int, xys ...plotter.XY) { p.InsertAll(i, xys) }
func (p *Plotter) InsertAll(i int, xys plotter.XYs) {
    if len(xys) == 0 {
        return
    }
    p.setXYs(append(p.xys[:i], append(xys, p.xys[i:]...)...))
}

func (p *Plotter) ReplaceAll(xys plotter.XYs) {
    p.setXYs(xys)
}

func (p *Plotter) setXYs(xys plotter.XYs) {
    p.mut.Lock()
    p.xys = xys
    defer p.mut.Unlock()
    
    select {
    case p.chChange <- time.Now():
    default:
    }
}

// PaddingX returns the padding ratio for the X axis.
func (p *Plotter) PaddingX() float64 { return p.padRatioX }

// PaddingY returns the padding ratio for the Y axis.
func (p *Plotter) PaddingY() float64 { return p.padRatioY }

// SetPaddingX sets the padding ratio for the X axis.
// Must be between -1 and 1;
func (p *Plotter) SetPaddingX(padding float64) error { return p.setPadding(padding, p.padRatioY) }

// SetPaddingX sets the data padding for the Y axis.
// Must be between -1 and 1;
func (p *Plotter) SetPaddingY(padding float64) error { return p.setPadding(p.padRatioX, padding) }

func (p *Plotter) setPadding(x, y float64) error {
    if x < -1 || x > 1 || y < -1 || y > 1 {
        return constants.ErrPaddingOutOfRange
    }
    p.padRatioX, p.padRatioY = x, y
    return nil
}

// Plot implements the plot.Plotter interface
func (p *Plotter) Plot(c draw.Canvas, plt *plot.Plot) {
    var xys plotter.XYs
    
    p.mut.Lock()
    if p.DataLimit <= 0 || p.xys.Len() < p.DataLimit {
        xys = p.xys[:]
    } else {
        xys = p.xys[p.xys.Len()-p.DataLimit:]
    }
    p.mut.Unlock()
    
    line, points, err := plotter.NewLinePoints(xys)
    utils.CheckErr(err)
    
    line.LineStyle = p.LineStyle
    points.GlyphStyleFunc = p.PointStyleFunc
    
    xmin, xmax, ymin, ymax := utils.FindMinMax(xys...)
    xpad, ypad := p.padRatioX*(xmax-xmin), p.padRatioY*(ymax-ymin)
    plt.X.Min, plt.X.Max = xmin-xpad, xmax+xpad
    plt.Y.Min, plt.Y.Max = ymin-ypad, ymax+ypad
    
    
    line.Plot(c, plt)
    points.Plot(c, plt)
}
