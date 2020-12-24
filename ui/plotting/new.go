package plotting

import (
    "gonum.org/v1/plot/plotter"
    "gonum.org/v1/plot/plotutil"
    "gonum.org/v1/plot/vg"
    "gonum.org/v1/plot/vg/draw"
    "image/color"
    "time"
)

// TODO Documentation

const (
	DefaultPaddingRatio = 0.03
)

var (
    DefaultLineStyle = draw.LineStyle{
        Color:    plotutil.Color(0),
        Width:    vg.Points(1),
        Dashes:   nil,
        DashOffs: 0,
    }
    DefaultPointStyle = draw.GlyphStyle{
        Color:  color.Black,
        Radius: vg.Points(1),
        Shape:  draw.CircleGlyph{},
    }
)

func NewPlotter(opts ...Option) (p *Plotter, err error) {
    cfg := config{
        lineStyle:  DefaultLineStyle,
        pointStyle: DefaultPointStyle,
        padX:       DefaultPaddingRatio,
        padY:       DefaultPaddingRatio,
    }
    
    for _, opt := range opts {
        opt(&cfg)
    }
    
    p = &Plotter{
        Name:           cfg.name,
        LineStyle:      cfg.lineStyle,
        PointStyleFunc: func(_ int) draw.GlyphStyle { return cfg.pointStyle },
        DataLimit:      cfg.limit,
        chChange:       make(chan time.Time, 1),
    }
    
    if err = p.SetPaddingX(cfg.padX); err != nil {
        return nil, err
    }
    if err = p.SetPaddingY(cfg.padY); err != nil {
        return nil, err
    }
    
    p.AppendAll(cfg.xys)
    
    return p, nil
}

type config struct {
    name       string
    xys        plotter.XYs
    limit      int
    lineStyle  draw.LineStyle
    pointStyle draw.GlyphStyle
    padX, padY float64
}

type Option func(opts *config)

func WithPaddingRatio(padX, padY float64) Option {
    return func(opts *config) { opts.padX, opts.padY = padX, padY }
}
func WithDataLimit(limit int) Option            { return func(opts *config) { opts.limit = limit } }
func WithLegend(name string) Option             { return func(opts *config) { opts.name = name } }
func WithData(xys ...plotter.XY) Option         { return func(opts *config) { opts.xys = xys } }
func WithLineStyle(style draw.LineStyle) Option { return func(opts *config) { opts.lineStyle = style } }
func WithPointStyle(style draw.GlyphStyle) Option {
    return func(opts *config) { opts.pointStyle = style }
}
func WithColorIdx(i int) Option { return WithColor(plotutil.Color(i)) }
func WithColor(c color.Color) Option {
    if c == nil {
        c = plotter.DefaultLineStyle.Color
    }
    return func(opts *config) {
        opts.lineStyle.Color = c
        opts.pointStyle.Color = c
    }
}
func WithDashesIdx(i int) Option { return WithDashes(plotutil.Dashes(i)) }
func WithDashes(dashes []vg.Length) Option {
    return func(opts *config) { opts.lineStyle.Dashes = dashes }
}
func WithShapeIdx(i int) Option { return WithShape(plotutil.Shape(i)) }
func WithShape(shape draw.GlyphDrawer) Option {
    return func(opts *config) { opts.pointStyle.Shape = shape }
}
func WithStyleIdx(i int) Option {
    return func(opts *config) {
        WithColorIdx(i)(opts)
        WithDashesIdx(i)(opts)
        WithShapeIdx(i)(opts)
    }
}
