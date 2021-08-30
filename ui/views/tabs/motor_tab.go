package tabs

import (
    "fmt"
    "gioui.org/layout"
    "gioui.org/op"
    "gioui.org/text"
    "gioui.org/unit"
    "gioui.org/widget"
    "gioui.org/widget/material"
    "github.com/ul-gaul/go-basestation/controller"
    "github.com/ul-gaul/go-basestation/data/packet"
    "github.com/ul-gaul/go-basestation/data/packet/state"
    "github.com/ul-gaul/go-basestation/ui/components/graph"
    "github.com/ul-gaul/go-basestation/ui/types"
    "github.com/ul-gaul/go-basestation/ui/vars"
    "image/color"
    "time"
)

var _ types.IComponent = (*MotorTab)(nil)

type MotorTab struct {
    last           packet.RocketPacket
    btns           map[state.Item]*widget.Clickable
    loading        map[state.Item]bool // temporary TODO remove and use state from last/current (RocketPacket)
    pressuresGraph *graph.Graph
}

func NewMotorTab() *MotorTab {
    return &MotorTab{
        btns:    make(map[state.Item]*widget.Clickable),
        loading: make(map[state.Item]bool),
        pressuresGraph: graph.New(
            graph.NewMultiSeriesRenderer(func(pkt packet.RocketPacket) (float64, []float64) {
                return float64(pkt.Time), []float64{
                    float64(pkt.Voltages.Item3.PSI()),
                    float64(pkt.Voltages.Item17.PSI()),
                    float64(pkt.Voltages.Item18.PSI()),
                    float64(pkt.Voltages.Item19.PSI()),
                }
            }, "Item 3", "Item 17", "Item 18", "Item 19"),
            graph.WithXAxis("Time"),
            graph.WithYAxis("Pressure (PSI)"),
            graph.WithTitle("Pressures"),
            graph.WithFixedRange(0, 1000),
            graph.WithXTicker(graph.TimeTicker),
            graph.WithLegend()),
    }
}

func (m *MotorTab) Tick(delta time.Duration) {
    m.last = vars.DataRoller.Last()
    m.pressuresGraph.Tick(delta)
}

func (m *MotorTab) Draw(gtx layout.Context) layout.Dimensions {
    return layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
        return layout.Flex{
            Axis:      layout.Horizontal,
            Spacing:   layout.SpaceBetween,
            Alignment: layout.Middle,
        }.Layout(gtx,
            // Col 1
            layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
                return layout.Flex{
                    Axis:      layout.Vertical,
                    Spacing:   layout.SpaceBetween,
                    Alignment: layout.Middle,
                }.Layout(gtx,
                    // Col 1 // Row 1
                    layout.Rigid(m.topLeft),
                    // Col 1 // Row 2
                    layout.Flexed(1, m.bottomLeft))
            }),
            
            // Col 2
            layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
                return layout.Flex{
                    Axis:      layout.Vertical,
                    Spacing:   layout.SpaceBetween,
                    Alignment: layout.Middle,
                }.Layout(gtx,
                    // Col 2 // Row 1
                    layout.Flexed(1, m.topRight),
                    // Col 2 // Row 2
                    layout.Flexed(1, m.bottomRight))
            }))
    })
}

/*****************************************************************************/
var d = time.Now()
func (m *MotorTab) topLeft(gtx layout.Context) layout.Dimensions {
    return layout.Flex{
        Axis:      layout.Horizontal,
        Spacing:   layout.SpaceBetween,
        Alignment: layout.Middle,
    }.Layout(gtx,
        layout.Flexed(1, m.flightDataSection),
        
        layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
            return layout.Dimensions{} // TODO or leave empty?
        }))
}

func (m *MotorTab) topRight(gtx layout.Context) layout.Dimensions {
    return m.pressuresGraph.Draw(gtx)
}

func (m *MotorTab) bottomLeft(gtx layout.Context) layout.Dimensions {
    return layout.Dimensions{} // TODO or leave empty?
}

func (m *MotorTab) bottomRight(gtx layout.Context) layout.Dimensions {
    return layout.Flex{
        Axis:      layout.Vertical,
        Spacing:   layout.SpaceEvenly,
        Alignment: layout.End,
    }.Layout(gtx,
        // For spacing
        layout.Flexed(0.3, blank),
        
        layout.Rigid(m.actuatorsSection),
        
        layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
            return layout.Flex{
                Axis:      layout.Horizontal,
                Spacing:   layout.SpaceBetween,
                Alignment: layout.Middle,
            }.Layout(gtx,
                layout.Flexed(1, m.pressuresSection),
                layout.Flexed(0.2, blank),
                layout.Flexed(1, m.voltagesSection),
                layout.Flexed(0.1, blank),
            )
        }),
        layout.Flexed(0.2, blank),
    )
}

/*****************************************************************************/

func blank(_ layout.Context) layout.Dimensions { return layout.Dimensions{} }

func headerLayout(str string) layout.Widget {
    header := material.Body1(vars.Theme, str)
    header.Font.Typeface = "Go"
    header.Font.Weight = text.Bold
    header.TextSize = header.TextSize.Scale(1.2)
    return header.Layout
}

func labelLayout(lbl string) layout.Widget {
    return material.Body1(vars.Theme, lbl).Layout
}

func valueLayout(val string) layout.Widget {
    v := material.Body1(vars.Theme, val)
    v.Font.Variant = "Mono"
    v.Alignment = text.End
    return func(gtx layout.Context) layout.Dimensions {
        op.InvalidateOp{}.Add(gtx.Ops)
        return v.Layout(gtx)
    }
}

func line(lbl, val string) layout.FlexChild {
    return layout.Rigid(func(gtx layout.Context) layout.Dimensions {
        return layout.Flex{
            Axis:      layout.Horizontal,
            Spacing:   layout.SpaceBetween,
            Alignment: layout.Middle,
        }.Layout(gtx,
            layout.Rigid(labelLayout(lbl)),
            layout.Rigid(valueLayout(val)))
    })
}

/*****************************************************************************/

func (m *MotorTab) flightDataSection(gtx layout.Context) layout.Dimensions {
    spacer := layout.Rigid(layout.Spacer{Height: unit.Dp(5)}.Layout)
    return layout.Flex{
        Axis:      layout.Vertical,
        Spacing:   layout.SpaceEnd,
        Alignment: layout.Start,
    }.Layout(gtx,
        layout.Rigid(headerLayout("Flight Data")),
        spacer,
        line("Time :", (time.Duration(m.last.Time)*time.Millisecond).Round(time.Millisecond*100).String()),
        spacer,
        line("Temperature :", fmt.Sprintf("%.1f °C", m.last.Temperature)),
        spacer,
        line("Humidity :", fmt.Sprintf("%.1f  %%", m.last.Humidity)),
        spacer,
        line("Atmos. Pressure :", fmt.Sprintf("%d Pa", m.last.AtmosPressure)),
        spacer,
        line("Altitude :", fmt.Sprintf("%.1f  m", m.last.Altitude)),
        spacer,
        line("Latitude :", fmt.Sprintf("%.6f", m.last.Latitude)),
        spacer,
        line("Longitude :", fmt.Sprintf("%.6f", m.last.Longitude)),
        spacer,
        line("GPS Satellites :", fmt.Sprintf("%d", m.last.Satellites)),
    )
}

func (m *MotorTab) actuatorsSection(gtx layout.Context) layout.Dimensions {
    return layout.Flex{
        Axis:      layout.Vertical,
        Spacing:   layout.SpaceSides,
        Alignment: layout.Start,
    }.Layout(gtx,
        layout.Rigid(headerLayout("Actuators")),
        layout.Rigid(func(gtx layout.Context) layout.Dimensions {
            return layout.Flex{
                Axis:      layout.Horizontal,
                Spacing:   layout.SpaceBetween,
                Alignment: layout.Middle,
            }.Layout(gtx,
                layout.Flexed(1, m.actuatorButton("A", state.I1)),
                layout.Flexed(1, m.actuatorButton("B", state.I2)),
                layout.Flexed(1, m.actuatorButton("C", state.I3)),
                layout.Flexed(1, m.actuatorButton("D", state.I4)),
                layout.Flexed(1, m.actuatorButton("E", state.I5)),
                layout.Flexed(1, m.actuatorButton("F", state.I6)),
                // layout.Flexed(1, m.actuatorButton("G", state.I7)),
                // layout.Flexed(1, m.actuatorButton("H", state.I8)),
            )
        }),
    )
}

func (m *MotorTab) actuatorButton(name string, item state.Item) layout.Widget {
    status := m.last.State.StatusOf(item)
    
    if _, ok := m.btns[item]; !ok {
        m.btns[item] = &widget.Clickable{}
    }
    
    if m.btns[item].Clicked() && !m.loading[item] {
        m.loading[item] = true
        controller.SetActuator(item, !status, func() {
            m.loading[item] = false
        })
    }
    
    btnStyle := material.ButtonLayout(vars.Theme, m.btns[item])
    txtStyle := material.Body1(vars.Theme, name)
    btnStyle.Background = color.NRGBA{R: 0xFF, A: 0xFF}
    txtStyle.Color = color.NRGBAModel.Convert(color.White).(color.NRGBA)
    txtStyle.Font.Weight = text.Bold
    if status {
        btnStyle.Background = color.NRGBA{G: 0xFF, A: 0xFF}
        txtStyle.Color = color.NRGBA{A: 0xFF}
    }
    
    return func(gtx layout.Context) layout.Dimensions {
        return layout.UniformInset(unit.Dp(4)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
            return btnStyle.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
                if m.loading[item] {
                    return layout.Flex{
                        Axis:      layout.Horizontal,
                        Spacing:   layout.SpaceAround,
                        Alignment: layout.Middle,
                    }.Layout(gtx,
                        layout.Rigid(func(gtx layout.Context) layout.Dimensions {
                            return layout.UniformInset(unit.Dp(7)).Layout(gtx, txtStyle.Layout)
                        }),
                        layout.Rigid(layout.Spacer{Width: unit.Dp(3)}.Layout),
                        layout.Rigid(material.Loader(vars.Theme).Layout))
                }
                return layout.UniformInset(unit.Dp(7)).Layout(gtx, txtStyle.Layout)
            })
        })
    }
}

func (m *MotorTab) pressuresSection(gtx layout.Context) layout.Dimensions {
    v := m.last.Voltages
    return layout.Flex{
        Axis:      layout.Vertical,
        Spacing:   layout.SpaceEvenly,
        Alignment: layout.Start,
    }.Layout(gtx,
        layout.Rigid(headerLayout("Pressures")),
        line("Item 3 :", fmt.Sprintf("%.1f PSI", v.Item3.PSI())),
        line("Item 17 :", fmt.Sprintf("%.1f PSI", v.Item17.PSI())),
        line("Item 18 :", fmt.Sprintf("%.1f PSI", v.Item18.PSI())),
        line("Item 19 :", fmt.Sprintf("%.1f PSI", v.Item19.PSI())),
    )
}

func (m *MotorTab) voltagesSection(gtx layout.Context) layout.Dimensions {
    v := m.last.Voltages
    return layout.Flex{
        Axis:      layout.Vertical,
        Spacing:   layout.SpaceEvenly,
        Alignment: layout.Start,
    }.Layout(gtx,
        layout.Rigid(headerLayout("Voltages")),
        line("Rocket :", fmt.Sprintf("%.2f V", v.Rocket)),
        line("Tower :", fmt.Sprintf("%.2f V", v.Tower)),
        line("Ignition Box :", fmt.Sprintf("%.2f V", v.IgnitionBox)),
    )
}
