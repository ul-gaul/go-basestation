package tabs

import (
    "gioui.org/layout"
    "gioui.org/op"
    "gioui.org/unit"
    "github.com/ul-gaul/go-basestation/data/packet"
    "github.com/ul-gaul/go-basestation/pool"
    "github.com/ul-gaul/go-basestation/ui/components/graph"
    "github.com/ul-gaul/go-basestation/ui/types"
    "reflect"
    "sync"
    "time"
)

type GraphID uint8

const (
    GraphAltitude GraphID = 0b0000_0001 << iota
    GraphCoords
    GraphTemperature
    GraphPressure
)

var _ types.IComponent = (*GeneralTab)(nil)

type GeneralTab struct {
    graphs map[GraphID]*graph.Graph
}

func Extract(xKey, yKey string) graph.Renderer { // TODO temporary func
    return graph.NewSimpleRenderer(func(pkt packet.RocketPacket) (float64, float64) {
        val := reflect.ValueOf(pkt)
        return val.FieldByName(xKey).Convert(reflect.TypeOf(0.0)).Float(),
            val.FieldByName(yKey).Convert(reflect.TypeOf(0.0)).Float()
    })
}

func NewGeneralTab() *GeneralTab {
    tab := &GeneralTab{
        graphs: make(map[GraphID]*graph.Graph),
    }
    
    tab.graphs[GraphAltitude] = graph.New(Extract("Time", "Altitude"),
        graph.WithXAxis("Time"),
        graph.WithYAxis("Altitude (m)"),
        graph.WithTitle("Altitude"),
        graph.WithFixedRange(0, 35_000),
        graph.WithXTicker(graph.TimeTicker))
    g := tab.graphs[GraphAltitude]
    _ = g
    
    tab.graphs[GraphTemperature] = graph.New(Extract("Time", "Temperature"),
        graph.WithXAxis("Time"),
        graph.WithYAxis("Temperature (°C)"),
        graph.WithTitle("Temperature"),
        graph.WithFixedRange(-30, 45),
        graph.WithXTicker(graph.TimeTicker))
    g = tab.graphs[GraphTemperature]
    
    tab.graphs[GraphPressure] = graph.New(Extract("Time", "AtmosPressure"),
        graph.WithXAxis("Time"),
        graph.WithYAxis("Pressure (Pa)"),
        graph.WithTitle("Atmospheric Pressure"),
        graph.WithFixedRange(0, 101325),
        graph.WithXTicker(graph.TimeTicker))
    g = tab.graphs[GraphPressure]
    
    tab.graphs[GraphCoords] = graph.New(Extract("Latitude", "Longitude"),
        graph.WithXAxis("Latitude"),
        graph.WithYAxis("Longitude"),
        graph.WithTitle("Coordinates"))
    
    return tab
}


func (tab *GeneralTab) Tick(delta time.Duration) {
    var wg sync.WaitGroup
    for _, g := range tab.graphs {
        wg.Add(1)
        gr := g
        err := pool.Frontend.Submit(func() {
            defer wg.Done()
            gr.Tick(delta)
        })
        if err != nil {
            wg.Done()
        }
    }
    wg.Wait()
}

func (tab *GeneralTab) Draw(gtx layout.Context) layout.Dimensions {
    // fonction qui retourne le layout du graphique spécifié
    flexedChart := func(id GraphID) layout.FlexChild {
        return layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
            op.InvalidateOp{}.Add(gtx.Ops)
            return layout.UniformInset(unit.Px(10)).Layout(gtx, tab.graphs[id].Draw)
        })
    }
    
    return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
        
        // Column 1
        layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
            return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
                // Col 1, Row 1
                flexedChart(GraphAltitude),
                // Col 1, Row 2
                flexedChart(GraphCoords),
            )
        }),
        
        // Column 2
        layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
            return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
                // Col 2, Row 1
                flexedChart(GraphTemperature),
                // Col 2, Row 2
                flexedChart(GraphPressure),
            )
        }),
    )
}
