package graph

import (
    "github.com/ul-gaul/go-basestation/data/packet"
    "github.com/wcharczuk/go-chart/v2"
    "math"
    "time"
)


type Ticker func([]packet.RocketPacket) []chart.Tick

func TimeTicker(packets []packet.RocketPacket) []chart.Tick {
    if len(packets) == 0 {
        return nil
    }
    
    min, max := math.Inf(+1), math.Inf(-1)
    for _, pkt := range packets {
        if min > float64(pkt.Time) {
            min = float64(pkt.Time)
        }
        
        if max < float64(pkt.Time) {
            max = float64(pkt.Time)
        }
    }
    
    const minor, major = 500.0, 1000.0
    rMin := math.Ceil(min/minor) * minor
    
    ticks := make([]chart.Tick, int(math.Ceil((max-rMin)/minor))+2)
    ticks[0], ticks[len(ticks)-1] = chart.Tick{Value: min}, chart.Tick{Value: max}
    lblMax := 0
    for i := 1; i < len(ticks)-1; i++ {
        val := float64(i-1)*minor + rMin
        ticks[i] = chart.Tick{Value: val}
        if i != 1 && math.Mod(val, major) == 0 {
            lbl := (time.Duration(val) * time.Millisecond).Round(minor * time.Millisecond).String()
            ticks[i].Label = lbl
            if len(lbl) > lblMax {
                lblMax = len(lbl)
            }
        }
    }
    
    return ticks
}