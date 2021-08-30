package graph

import (
	"github.com/ul-gaul/go-basestation/data/packet"
	"github.com/ul-gaul/go-basestation/ui/vars"
	"github.com/wcharczuk/go-chart/v2"
	"github.com/wcharczuk/go-chart/v2/drawing"
)

type Renderer func([]packet.RocketPacket) []chart.Series

type SimpleAdapter func(packet.RocketPacket) (float64, float64)

func NewSimpleRenderer(adapter SimpleAdapter) Renderer {
	return func(packets []packet.RocketPacket) []chart.Series {
		xValues := make([]float64, len(packets))
		yValues := make([]float64, len(packets))

		for i, pkt := range packets {
			xValues[i], yValues[i] = adapter(pkt)
		}

		mainSeries := chart.ContinuousSeries{
			Style: chart.Style{
				StrokeColor: chart.GetDefaultColor(0),
			},
			XValues: xValues,
			YValues: yValues,
		}

		return []chart.Series{
			mainSeries,
		}
	}
}

type MultiSeriesAdapter func(packet.RocketPacket) (x float64, yValues []float64)

func NewMultiSeriesRenderer(adapter MultiSeriesAdapter, names ...string) Renderer {
	return func(packets []packet.RocketPacket) []chart.Series {
		var yValuesList [][]float64
		xValues := make([]float64, len(packets))

		for i, pkt := range packets {
			var yValues []float64
			xValues[i], yValues = adapter(pkt)

			for j, yValue := range yValues {
				if j >= len(yValuesList) {
					yValuesList = append(yValuesList, make([]float64, len(packets)))
				}
				yValuesList[j][i] = yValue
			}
		}

		series := make([]chart.Series, len(yValuesList))
		for i, yValues := range yValuesList {
			r, g, b, a := vars.Color(i).RGBA()
			s := chart.ContinuousSeries{
				Style: chart.Style{
					StrokeColor: drawing.ColorFromAlphaMixedRGBA(r, g, b, a),
				},
				XValues: xValues,
				YValues: yValues,
			}
			
			if i < len(names) {
				s.Name = names[i]
			}
			
			series[i] = s
		}

		return series
	}
}
