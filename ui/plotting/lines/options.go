package lines

import (
    "gonum.org/v1/plot/plotutil"
    "gonum.org/v1/plot/vg"
    "gonum.org/v1/plot/vg/draw"
    "image/color"
)

type LineOption func(line draw.LineStyle)

func LineWidth(w vg.Length) LineOption {
    return func(line draw.LineStyle) {
        line.Width = w
    }
}

func LineColorIdx(i int) LineOption { return LineColor(plotutil.Color(i)) }
func LineColor(c color.Color) LineOption {
    return func(line draw.LineStyle) {
        line.Color = c
    }
}

func LineDashesIdx(i int) LineOption { return LineDashes(plotutil.Dashes(i)...) }
func LineDashes(dashes ...vg.Length) LineOption {
    return func(line draw.LineStyle) {
        line.Dashes = dashes
    }
}

func LineDashOffs(doffs vg.Length) LineOption {
    return func(line draw.LineStyle) {
        line.DashOffs = doffs
    }
}
