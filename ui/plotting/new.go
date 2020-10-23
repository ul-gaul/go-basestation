package plotting

import (
    "gonum.org/v1/plot/plotter"
    "gonum.org/v1/plot/plotutil"
    "gonum.org/v1/plot/vg"
    "gonum.org/v1/plot/vg/draw"
    "image/color"
)

// TODO Documentation

func NewPlotter(opts ...Option) (p *Plotter, err error) {
    cfg := config{
        lineStyle:  plotter.DefaultLineStyle,
        pointStyle: plotter.DefaultGlyphStyle,
    }
    
    for _, opt := range opts {
        opt(&cfg)
    }
    
    p = &Plotter{name: cfg.name}
    p.line, p.points, err = plotter.NewLinePoints(new(plotter.XYs))
    if err != nil {
        return nil, err
    }
    p.line.LineStyle = cfg.lineStyle
    p.points.GlyphStyle = cfg.pointStyle
    p.AppendAll(cfg.xys)
    
    return p, nil
}

type config struct {
    name       string
    xys        plotter.XYs
    lineStyle  draw.LineStyle
    pointStyle draw.GlyphStyle
}

type Option func(opts *config)

func WithLegend(name string) Option             { return func(opts *config) { opts.name = name } }
func WithData(xys ...plotter.XY) Option         { return func(opts *config) { opts.xys = xys } }
func WithLineStyle(style draw.LineStyle) Option { return func(opts *config) { opts.lineStyle = style } }
func WithPointStyle(style draw.GlyphStyle) Option {
    return func(opts *config) { opts.pointStyle = style }
}
func WithColorIdx(i int) Option                 { return WithColor(plotutil.Color(i)) }
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
