package vgg

import (
    "image"
    "image/color"
    
    "gioui.org/f32"
    "gioui.org/layout"
    "gioui.org/op/paint"
    
    "gonum.org/v1/plot/vg"
    "gonum.org/v1/plot/vg/vgimg"
)

// Canvas implements the vg.Canvas interface,
// drawing to an image.Image using vgimg and painting that image
// into a Gioui context.
type Canvas struct {
    *vgimg.Canvas
}

// NewCanvas returns a new image canvas with the provided dimensions and options.
// The currently accepted options are UseDPI and UseBackgroundColor.
// If the resolution or background color are not specified, defaults are used.
func NewCanvas(opts ...option) *Canvas {
    cfg := &config{
        dpi: vgimg.DefaultDPI,
        bkg: color.White,
        w: vgimg.DefaultWidth,
        h: vgimg.DefaultHeight,
    }
    for _, opt := range opts {
        opt(cfg)
    }
    c := &Canvas{
        Canvas: vgimg.NewWith(
            vgimg.UseDPI(cfg.dpi),
            vgimg.UseWH(cfg.w, cfg.h),
            vgimg.UseBackgroundColor(cfg.bkg),
        ),
    }
    return c
}

type config struct {
    dpi int
    bkg color.Color
    w, h  vg.Length
}

type option func(*config)



// UseWH specifies the width and height of the canvas.
// The size is rounded up to the nearest pixel.
func UseWH(w, h vg.Length) option {
    if w <= 0 || h <= 0 {
        panic("w and h must both be > 0.")
    }
    return func(c *config) {
        c.w = w
        c.h = h
    }
}

// UseDPI sets the dots per inch of a canvas. It should only be
// used as an option argument when initializing a new canvas.
func UseDPI(dpi int) option {
    if dpi <= 0 {
        panic("DPI must be > 0.")
    }
    return func(c *config) {
        c.dpi = dpi
    }
}

// UseBackgroundColor specifies the image background color.
// Without UseBackgroundColor, the default color is white.
func UseBackgroundColor(c color.Color) option {
    return func(cfg *config) {
        cfg.bkg = c
    }
}

func (c *Canvas) pt32(p vg.Point) f32.Point {
    _, h := c.Size()
    dpi := c.DPI()
    return f32.Point{
        X: float32(p.X.Dots(dpi)),
        Y: float32(h.Dots(dpi) - p.Y.Dots(dpi)),
    }
}

// Paint paints the canvas' content on the screen.
func (c *Canvas) Layout(gtx layout.Context) layout.Dimensions {
    r32 := f32.Rect(
        float32(gtx.Constraints.Min.X),
        float32(gtx.Constraints.Min.Y),
        float32(gtx.Constraints.Max.X),
        float32(gtx.Constraints.Max.Y))
    
    paint.NewImageOp(c.Canvas.Image()).Add(gtx.Ops)
    paint.PaintOp{Rect: r32}.Add(gtx.Ops)
    
    return layout.Dimensions{
        Size: image.Pt(int(r32.Max.X), int(r32.Max.Y)),
    }
}
