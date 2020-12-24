package utils

import (
    "gonum.org/v1/plot/plotter"
)

func FindMinMax(xys ...plotter.XY) (xmin, xmax, ymin, ymax float64) {
    for i, xy := range xys {
        if i == 0 {
            xmin = xy.X
            xmax = xy.X
            ymin = xy.Y
            ymax = xy.Y
        } else {
            if xy.X > xmax {
                xmax = xy.X
            } else if xy.X < xmin {
                xmin = xy.X
            }
    
            if xy.Y > ymax {
                ymax = xy.Y
            } else if xy.Y < ymin {
                ymin = xy.Y
            }
        }
    }
    return xmin, xmax, ymin, ymax
}
