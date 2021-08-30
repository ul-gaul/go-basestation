package graph

import (
	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/widget"
	"github.com/ul-gaul/go-basestation/cfg"
	"github.com/ul-gaul/go-basestation/data/packet"
	"github.com/ul-gaul/go-basestation/ui/types"
	"github.com/ul-gaul/go-basestation/ui/vars"
	"github.com/wcharczuk/go-chart/v2"
	"image"
	"time"
)

// TODO documentation

var _ types.IComponent = (*Graph)(nil)

type Graph struct {
	*chart.Chart
	Renderer
	XTicker, YTicker Ticker
	imgOp paint.ImageOp
}

func New(renderer Renderer, options ...func(*Graph)) *Graph {
	h := float64(cfg.Frontend.Graph.BaseWidth)
	h *= float64(chart.DefaultChartHeight) / float64(chart.DefaultChartWidth)
	
	g := &Graph{
		Chart: &chart.Chart{
			Width:  cfg.Frontend.Graph.BaseWidth,
			Height: int(h),
		},
		Renderer: renderer,
	}
	g.imgOp = paint.NewImageOp(image.NewRGBA(
		image.Rect(0, 0, g.GetWidth(), g.GetHeight())))

	for _, opt := range options {
		opt(g)
	}

	return g
}

func (g *Graph) Image() image.Image {
	var packets []packet.RocketPacket
	vars.DataRoller.Do(func(pkt packet.RocketPacket) {
		packets = append(packets, pkt)
	})
	
	if g.XTicker != nil {
		g.XAxis.Ticks = g.XTicker(packets)
		g.XAxis.GridLines = createGridLines(g.XAxis.Ticks,
			g.XAxis.GridMajorStyle, g.XAxis.GridMinorStyle)
	}
	
	if g.YTicker != nil {
		g.YAxis.Ticks = g.YTicker(packets)
		g.YAxis.GridLines = createGridLines(g.YAxis.Ticks,
			g.YAxis.GridMajorStyle, g.YAxis.GridMinorStyle)
	}
	
	g.Series = g.Renderer(packets)

	iw := &chart.ImageWriter{}
	_ = g.Render(chart.PNG, iw)
	m, _ := iw.Image()

	return m
}

func createGridLines(ticks []chart.Tick, major, minor chart.Style) []chart.GridLine {
	if ticks == nil || len(ticks) == 0 {
		return nil
	}
	
	gridLines := make([]chart.GridLine, len(ticks))
	for i, tick := range ticks {
		gridLines[i] = chart.GridLine{
			IsMinor: len(tick.Label) == 0,
			Style: minor,
			Value:   tick.Value,
		}
		
		if len(tick.Label) > 0 {
			gridLines[i].Style = major
		}
	}
	
	return gridLines
}

func (g *Graph) Tick(_ time.Duration) {
	g.imgOp = paint.NewImageOp(g.Image())
}

func (g *Graph) Draw(gtx layout.Context) layout.Dimensions {
	{ // Update width & height
		dRatio := float64(gtx.Constraints.Max.Y) / float64(gtx.Constraints.Max.X)
		w := float64(cfg.Frontend.Graph.BaseWidth)
		h := dRatio * w
		
		if int(w) != g.GetWidth() || int(h) != g.GetHeight() {
			g.Width, g.Height = int(w), int(h)
		}
	}
	
	{ // Update DPI
		dpi := float64(gtx.Metric.PxPerDp) * chart.DefaultDPI * cfg.Frontend.Graph.Scale
		if dpi != g.GetDPI() {
			g.DPI = dpi
		}
	}
	
	return widget.Image{
		Src:      g.imgOp,
		Fit:      widget.Fill,
		Position: layout.Center,
		Scale:    float32(160.0 / g.GetDPI()),
	}.Layout(gtx)
}
