package ticker

import (
    . "gonum.org/v1/plot"
    "math"
    "strconv"
)

const (
    DefaultSuggestedTicks = 5
    DefaultContainment = WithinData
)

// Ticks is suitable for the Tick.Marker field of an Axis,
// it returns a reasonable default set of tick marks.
type Ticks struct{
    suggestedTicks int
    containment Containment
}

// var _ Ticker = (Ticks)(nil)
var _ Ticker = Ticks{}

func NewTicker(suggestedTicks int, containment Containment) Ticker {
    suggestedTicks = maxInt(0, suggestedTicks)
    if suggestedTicks == 0 {
        suggestedTicks = DefaultSuggestedTicks
    }
    
    if containment != Free && containment != ContainData && containment != WithinData {
        containment = DefaultContainment
    }
    
    return Ticks{suggestedTicks, containment}
}

// Ticks returns Ticks in the specified range.
func (t Ticks) Ticks(min, max float64) []Tick {
    if max <= min {
        panic("illegal range")
    }
    
    labels, step, q, mag := talbotLinHanrahan(min, max, t.suggestedTicks, t.containment, nil, nil, nil)
    majorDelta := step * math.Pow10(mag)
    if q == 0 {
        // Simple fall back was chosen, so
        // majorDelta is the label distance.
        majorDelta = labels[1] - labels[0]
    }
    
    // Choose a reasonable, but ad
    // hoc formatting for labels.
    fc := byte('f')
    var off int
    if mag < -1 || 6 < mag {
        off = 1
        fc = 'g'
    }
    if math.Trunc(q) != q {
        off += 2
    }
    prec := minInt(6, maxInt(off, -mag))
    ticks := make([]Tick, len(labels))
    for i, v := range labels {
        ticks[i] = Tick{Value: v, Label: strconv.FormatFloat(v, fc, prec, 64)}
    }
    
    var minorDelta float64
    // See talbotLinHanrahan for the values used here.
    switch step {
    case 1, 2.5:
        minorDelta = majorDelta / 5
    case 2, 3, 4, 5:
        minorDelta = majorDelta / step
    default:
        if majorDelta/2 < dlamchP {
            return ticks
        }
        minorDelta = majorDelta / 2
    }
    
    // Find the first minor tick not greater
    // than the lowest data value.
    var i float64
    for labels[0]+(i-1)*minorDelta > min {
        i--
    }
    // Add ticks at minorDelta intervals when
    // they are not within minorDelta/2 of a
    // labelled tick.
    for {
        val := labels[0] + i*minorDelta
        if val > max {
            break
        }
        found := false
        for _, t := range ticks {
            if math.Abs(t.Value-val) < minorDelta/2 {
                found = true
            }
        }
        if !found {
            ticks = append(ticks, Tick{Value: val})
        }
        i++
    }
    
    return ticks
}

func minInt(a, b int) int {
    if a < b {
        return a
    }
    return b
}

func maxInt(a, b int) int {
    if a > b {
        return a
    }
    return b
}