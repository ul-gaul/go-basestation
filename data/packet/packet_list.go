package packet

import (
    . "gonum.org/v1/plot/plotter"
    "time"
)

type PacketList []RocketPacket

func timeToDuration(t uint64) time.Duration {
    return time.Duration(t) * time.Millisecond
}

func (pl PacketList) AltitudeData() XYs {
    xys := make(XYs, len(pl))
    for i, pkt := range pl {
        xys[i] = XY{timeToDuration(pkt.Time).Seconds(), pkt.Altitude}
    }
    return xys
}

func (pl PacketList) PressureData() XYs {
    xys := make(XYs, len(pl))
    for i, pkt := range pl {
        xys[i] = XY{timeToDuration(pkt.Time).Seconds(), pkt.Pressure}
    }
    return xys
}

func (pl PacketList) PositionData() XYs {
    xys := make(XYs , len(pl))
    for i, pkt := range pl {
        xys[i] = XY{pkt.Longitude, pkt.Latitude}
    }
    return xys
}
