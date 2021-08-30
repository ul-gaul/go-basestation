package vars

import (
    "image/color"
)

// DefaultColors is a set of colors used by the Color function.
var DefaultColors = DarkColors

var DarkColors = []color.Color{
    rgb(238, 46, 47),
    rgb(0, 140, 72),
    rgb(24, 90, 169),
    rgb(244, 125, 35),
    rgb(102, 44, 145),
    rgb(162, 29, 33),
    rgb(180, 56, 148),
}

var SoftColors = []color.Color{
    rgb(241, 90, 96),
    rgb(122, 195, 106),
    rgb(90, 155, 212),
    rgb(250, 167, 91),
    rgb(158, 103, 171),
    rgb(206, 112, 88),
    rgb(215, 127, 180),
}

func rgb(r, g, b uint8) color.RGBA {
    return color.RGBA{R: r, G: g, B: b, A: 255}
}

// Color returns the ith default color, wrapping
// if i is less than zero or greater than the max
// number of colors in the DefaultColors slice.
func Color(i int) color.Color {
    n := len(DefaultColors)
    if i < 0 {
        return DefaultColors[i%n+n]
    }
    return DefaultColors[i%n]
}
