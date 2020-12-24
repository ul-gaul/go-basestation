package plotting

import (
    "gonum.org/v1/plot"
    "gonum.org/v1/plot/plotter"
    "gonum.org/v1/plot/vg/draw"
    "math"
    "sync"
    "time"
    
    "github.com/ul-gaul/go-basestation/utils"
)

var _ plot.DataRanger = (*Plotter)(nil)
var _ plot.Plotter = (*Plotter)(nil)

type Plotter struct {
    name                   string
    line                   *plotter.Line
    points                 *plotter.Scatter
    xys                    plotter.XYs
    xmin, xmax, ymin, ymax float64
    chChange               chan time.Time
    mut                    sync.Mutex
}

func (p *Plotter) Name() string      { return p.name }
func (p *Plotter) Data() plotter.XYs { return p.xys }

func (p *Plotter) DataRange() (xmin, xmax, ymin, ymax float64) {
    return p.xmin, p.xmax, p.ymin, p.ymax
}

func (p *Plotter) Prepend(xys ...plotter.XY)       { p.PrependAll(xys) }
func (p *Plotter) PrependAll(xys plotter.XYs)      { p.InsertAll(0, xys) }
func (p *Plotter) Append(xys ...plotter.XY)        { p.AppendAll(xys) }
func (p *Plotter) AppendAll(xys plotter.XYs)       { p.InsertAll(len(p.xys), xys) }
func (p *Plotter) Insert(i int, xys ...plotter.XY) { p.InsertAll(i, xys) }
func (p *Plotter) InsertAll(i int, xys plotter.XYs) {
    if len(xys) == 0 {
        return
    }
    
    xmin, xmax, ymin, ymax := utils.FindMinMax(xys...)
    p.xmin = math.Min(p.xmin, xmin)
    p.xmax = math.Max(p.xmax, xmax)
    p.ymin = math.Min(p.ymin, ymin)
    p.ymax = math.Max(p.ymax, ymax)
    
    p.setXYs(append(p.xys[:i], append(xys, p.xys[i:]...)...))
}

func (p *Plotter) ReplaceAll(xys plotter.XYs) {
    p.xmin, p.xmax, p.ymin, p.ymax = utils.FindMinMax(xys...)
    p.setXYs(xys)
}

func (p *Plotter) setXYs(xys plotter.XYs) {
    p.mut.Lock()
    p.xys = xys
    p.line.XYs = xys
    p.points.XYs = xys
    defer p.mut.Unlock()
    
    select {
    case p.chChange <- time.Now():
    default:
    }
}

// Plot implements the plot.Plotter interface
func (p *Plotter) Plot(c draw.Canvas, plt *plot.Plot) {
    p.mut.Lock()
    line, points, err := plotter.NewLinePoints(p.xys)
    utils.CheckErr(err)
    p.mut.Unlock()
    
    line.LineStyle = p.LineStyle()
    points.GlyphStyle = p.PointStyle()
    
    line.Plot(c, plt)
    points.Plot(c, plt)
}

/*********************** STYLING ***********************/

func (p *Plotter) LineStyle() draw.LineStyle                 { return p.line.LineStyle }
func (p *Plotter) PointStyle() draw.GlyphStyle               { return p.points.GlyphStyle }
func (p *Plotter) PointStyleFunc() func(int) draw.GlyphStyle { return p.points.GlyphStyleFunc }
func (p *Plotter) SetPointStyleFunc(fnc func(int) draw.GlyphStyle) {
    p.points.GlyphStyleFunc = fnc
}
