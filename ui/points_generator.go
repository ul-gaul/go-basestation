package ui

import (
    "gonum.org/v1/plot/plotter"
    "math"
    "math/rand"
    
    "github.com/ul-gaul/go-basestation/ui/plotting"
)

// randomPoints returns some random x, y points.
func randomPointsNeg(n int) plotter.XYs {
    pts := make(plotter.XYs, n)
    last := rand.Float64()
    for i := range pts {
        last += math.Copysign(10*rand.Float64(), rand.Float64()-0.5)
        pts[i].X = last
        pts[i].Y = last + math.Copysign(10*rand.Float64(), rand.Float64()-0.5)
    }
    return pts
}

func addRandomPoints(p *plotting.Plotter, n int) {
    pts := make(plotter.XYs, n)
    var last plotter.XY
    data := p.Data()
    if len(data) > 0 { last = data[len(data)-1] }
    for i := range pts {
        last = plotter.XY{
            X: last.X + 1,
            Y: 10 * rand.Float64(),
        }
        pts[i] = last
    }
    p.Append(pts...)
}

// squarePlot *TEMP* génère des points
func squarePlot(i int) plotter.XYs {
    const b, d = 5.0, 2.0
    dist := b + (d * float64(i))
    xys := make(plotter.XYs, 6)
    xys[0] = plotter.XY{0, 0}
    
    var calcSign = func(v float64) float64 {
        // (-1) ** ( floor( (v%4) / 2 ) % 2 )
        return math.Mod(math.Pow(-1, math.Floor(math.Mod(v, 4)/2)), 2)
    }
    var calc = func(v float64) float64 {
        return calcSign(v) * dist
    }
    
    for j := 1; j < len(xys); j++ {
        xys[j].X = calc(float64(i) + float64(j) - 1)
        xys[j].Y = calc(float64(i) + float64(j))
    }
    
    last := &xys[len(xys)-1]
    if i%2 == 0 {
        last.Y += math.Copysign(1, last.Y)
    } else {
        last.X += math.Copysign(1, last.X)
    }
    
    return xys
}
