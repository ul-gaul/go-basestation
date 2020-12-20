package plotting

import (
    log "github.com/sirupsen/logrus"
    "gonum.org/v1/plot/vg/vgimg"
    "image/color"
    
    "gonum.org/v1/plot/vg"
)

type drawerCfg struct {
    dpi int
    bkg color.Color
    w, h  vg.Length
}

// newCanvas returns a new image canvas with the provided dimensions and options.
// The currently accepted options are UseDPI and UseBackgroundColor.
// If the resolution or background color are not specified, defaults are used.
func newCanvas(opts ...option) *vgimg.Canvas {
    cfg := &drawerCfg{
        dpi: vgimg.DefaultDPI,
        bkg: color.White,
        w: vgimg.DefaultWidth,
        h: vgimg.DefaultHeight,
    }
    for _, opt := range opts {
        opt(cfg)
    }
    return vgimg.NewWith(
        vgimg.UseDPI(cfg.dpi),
        vgimg.UseWH(cfg.w, cfg.h),
        vgimg.UseBackgroundColor(cfg.bkg),
    )
}

type option func(*drawerCfg)


// UseWH specifies the width and height of the canvas.
// The size is rounded up to the nearest pixel.
func UseWH(w, h vg.Length) option {
    if w <= 0 || h <= 0 {
        log.Panicln("w and h must both be > 0.")
    }
    return func(c *drawerCfg) {
        c.w = w
        c.h = h
    }
}

// UseDPI sets the dots per inch of a canvas. It should only be
// used as an option argument when initializing a new canvas.
func UseDPI(dpi int) option {
    if dpi <= 0 {
        log.Panicln("DPI must be > 0.")
    }
    return func(c *drawerCfg) {
        c.dpi = dpi
    }
}

// UseBackgroundColor specifies the image background color.
// Without UseBackgroundColor, the default color is white.
func UseBackgroundColor(c color.Color) option {
    return func(cfg *drawerCfg) {
        cfg.bkg = c
    }
}
