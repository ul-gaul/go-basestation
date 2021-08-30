package graph

import (
	"github.com/ul-gaul/go-basestation/ui/vars"
	"github.com/wcharczuk/go-chart/v2"
	"github.com/wcharczuk/go-chart/v2/drawing"
	"math"
	"strconv"
)

const formatterMaxPrecision = 3

func valueFormatter(v interface{}) string {
	f := v.(float64)
	f = math.Round(f * math.Pow10(formatterMaxPrecision)) / math.Pow10(formatterMaxPrecision)
	return strconv.FormatFloat(f, 'f', -1, 64)
}

func WithYAxis(axisName string) func(*Graph) {
	return func(g *Graph) {
		g.YAxis = chart.YAxis{
			Name: axisName,
			Style: chart.Style{
				FontSize: 7,
				TextWrap: chart.TextWrapNone,
			},
			GridMajorStyle: chart.Style{
				StrokeWidth: 1,
				StrokeColor: drawing.ColorFromAlphaMixedRGBA(0xAA, 0xAA, 0xAA, 112),
			},
			GridMinorStyle: chart.Style{
				StrokeWidth: 1,
				StrokeColor: drawing.ColorFromAlphaMixedRGBA(0xAA, 0xAA, 0xAA, 48),
			},
			ValueFormatter: valueFormatter,
		}
		g.YAxis.NameStyle.FontSize = 7
		g.YAxis.NameStyle.FontColor = drawing.ColorFromAlphaMixedRGBA(vars.Theme.Fg.RGBA()).WithAlpha(128)
	}
}

func WithXAxis(axisName string) func(*Graph) {
	return func(g *Graph) {
		g.XAxis = chart.XAxis{
			Name: axisName,
			Style: chart.Style{
				FontSize: 7,
			},
			GridMajorStyle: chart.Style{
				StrokeWidth: 1,
				StrokeColor: drawing.ColorFromAlphaMixedRGBA(0xAA, 0xAA, 0xAA, 112),
			},
			GridMinorStyle: chart.Style{
				StrokeWidth: 1,
				StrokeColor: drawing.ColorFromAlphaMixedRGBA(0xAA, 0xAA, 0xAA, 48),
			},
			ValueFormatter: valueFormatter,
		}
		g.XAxis.NameStyle.FontSize = 7
		g.XAxis.NameStyle.FontColor = drawing.ColorFromAlphaMixedRGBA(vars.Theme.Fg.RGBA()).WithAlpha(128)
	}
}

func WithTitle(title string) func(*Graph) {
	return func(g *Graph) {
		g.Title = title
		g.TitleStyle = chart.Style{
			FontSize: 10,
		}
	}
}

func WithFixedRange(min, max float64) func(*Graph) {
	return func(g *Graph) {
		g.YAxis.Range = &chart.ContinuousRange{
			Min: min,
			Max: max,
		}
	}
}


func WithXTicker(ticker Ticker) func(*Graph) {
	return func(g *Graph) {
		g.XTicker = ticker
	}
}

func WithYTicker(ticker Ticker) func(*Graph) {
	return func(g *Graph) {
		g.YTicker = ticker
	}
}

func WithLegend(styles ...chart.Style) func(*Graph) {
	var style chart.Style
	if len(styles) > 0 {
		style = styles[0]
	}
	
	style = style.InheritFrom(chart.Style{
		FontSize: 7,
		TextLineSpacing: 5,
		StrokeWidth: 0.5,
	})
	
	return func(g *Graph) {
		g.Elements = append(g.Elements, chart.Legend(g.Chart, style))
	}
}