package plotting

import (
    assertions "github.com/stretchr/testify/assert"
    "gonum.org/v1/plot/plotter"
    "gonum.org/v1/plot/plotutil"
    "image/color"
    "testing"
)

func TestWithColor(t *testing.T) {
    assert := assertions.New(t)
    type test struct {
        color color.Color
        opts  []Option
    }
    var newTest = func(c color.Color) test {
        return test{c, []Option{WithColor(c)}}
    }
    tests := []test{
        {color: plotter.DefaultLineStyle.Color},
        {color: plotter.DefaultGlyphStyle.Color},
        {plotutil.Color(1), []Option{WithColorIdx(1)}},
        newTest(color.White),
        newTest(color.Black),
        newTest(color.RGBA{255, 0, 0, 255}),
        newTest(color.RGBA{0, 255, 0, 255}),
        newTest(color.RGBA{0, 0, 255, 255}),
        newTest(color.RGBA{127, 127, 127, 127}),
        {plotter.DefaultLineStyle.Color, []Option{WithColor(nil)}},
    }
    
    for _, tt := range tests {
        p, err := NewPlotter(tt.opts...)
        assert.Nil(err)
        assert.Equal(tt.color, p.points.Color)
        assert.Equal(tt.color, p.line.Color)
    }
}

func TestWithColorIdx(t *testing.T) {
    assert := assertions.New(t)
    tests := []int{-9_999_999, -1, 0, nil, 1, 9_999_999}
    
    for _, tt := range tests {
        p, err := NewPlotter(WithColorIdx(tt))
        assert.Nil(err)
        assert.Equal(plotutil.Color(tt), p.points.Color)
        assert.Equal(plotutil.Color(tt), p.line.Color)
    }
}

func TestWithDashes(t *testing.T) {
    // TODO
}

func TestWithDashesIdx(t *testing.T) {
    // TODO
}

func TestWithLegend(t *testing.T) {
    // TODO
}

func TestWithLineStyle(t *testing.T) {
    // TODO
}

func TestWithData(t *testing.T) {
    // TODO test xys
    // TODO test min/max
}

func TestWithPointStyle(t *testing.T) {
    // TODO
}

func TestWithShape(t *testing.T) {
    // TODO
}

func TestWithShapeIdx(t *testing.T) {
    // TODO
}

func TestWithStyleIdx(t *testing.T) {
    // TODO
}
