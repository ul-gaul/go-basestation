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
    // DefaultLineStyle est le style par défaut des lignes
    DefaultLineStyle = draw.LineStyle{
        Color:    plotutil.Color(0),
        Width:    vg.Points(1),
        Dashes:   nil,
        DashOffs: 0,
    }
    // DefaultPointStyle est le style par défaut des points
    DefaultPointStyle = draw.GlyphStyle{
        Color:  color.Black,
        Radius: vg.Points(1),
        Shape:  draw.CircleGlyph{},
    }
)

// NewPlotter crée un nouveau Plotter avec les options spécifiées
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

// WithPaddingRatio spécifie une marge dans le graphique
func WithPaddingRatio(padX, padY float64) Option {
    return func(opts *config) { opts.padX, opts.padY = padX, padY }
}
// WithDataLimit limite le nombre de données affichées dans le graphique
func WithDataLimit(limit int) Option            { return func(opts *config) { opts.limit = limit } }
// WithLegend spécifie une légende pour les données de ce plotter
func WithLegend(name string) Option             { return func(opts *config) { opts.name = name } }
// WithData spécifie des données initiales
func WithData(xys ...plotter.XY) Option         { return func(opts *config) { opts.xys = xys } }
// WithLineStyle spécifie le style de ligne
func WithLineStyle(style draw.LineStyle) Option { return func(opts *config) { opts.lineStyle = style } }
// WithPointStyle spécifie le style des points
func WithPointStyle(style draw.GlyphStyle) Option {
    return func(opts *config) { opts.pointStyle = style }
}
// WithColorIdx spécifie la couleur des lignes et des points (voir plotutil.Color)
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

// WithDashesIdx spécifie le style de pointillés (voir plotutil.Dashes)
func WithDashesIdx(i int) Option { return WithDashes(plotutil.Dashes(i)) }
func WithDashes(dashes []vg.Length) Option {
    return func(opts *config) { opts.lineStyle.Dashes = dashes }
}
// WithShapeIdx spécifie la forme des points (voir plotutil.Shape)
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
