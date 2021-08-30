package controller

import (
	"errors"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/panjf2000/ants/v2"
	"github.com/ul-gaul/go-basestation/data/manager"
	"github.com/ul-gaul/go-basestation/data/packet"
	"github.com/ul-gaul/go-basestation/data/packet/state"
	"github.com/ul-gaul/go-basestation/utils"
	"math"
	"math/rand"
	"time"
)

// TODO documentation

var (
	genRunning bool
	rdmState state.State
)

func Generate() error {
	if CurrentMode() != MODE_NONE {
		return errors.New("already running")
	}

	genRunning = true
	utils.CheckErr(ants.Submit(genLoop))
	return nil
}

func IsGeneratorStarted() bool {
	return genRunning
}

func StopGenerator() {
	genRunning = false
}

func genLoop() {
	r := func(i uint64, offset, d, min, max float64) float64 {
		f := float64(i) / 1000.0
		f += math.Pi * 2 * offset
		diff := max - min
		return utils.Constrain(
			math.Sin(f)*diff/2.5+diff/2+gofakeit.Float64Range(-d, d),
			min, max)
	}

	offsets := make([]float64, 20)
	for i := range offsets {
		offsets[i] = rand.Float64()
	}

	var ms uint64
	ch := time.Tick(10 * time.Millisecond)
	for genRunning {
		select {
		case <-ch:
			manager.Data(packet.RocketPacket{
				Time:          packet.Milliseconds(ms),
				Altitude:      float32(r(ms/5, offsets[0], 3, 0, 30000)),
				Latitude:      float32(gofakeit.Latitude()),
				Longitude:     float32(gofakeit.Longitude()),
				Temperature:   float32(r(ms, offsets[1], .5, 0, 45)),
				AtmosPressure: int32(r(ms*2, offsets[2], 2000, 0, 101325)),
				Voltages: packet.Voltages{
					Rocket:      gofakeit.Float32Range(0, 15),
					Tower:       gofakeit.Float32Range(0, 15),
					IgnitionBox: gofakeit.Float32Range(0, 15),
					Item3:       packet.VoltagePSI(r(ms, 0.0/4, 0.2, 0, 5)),
					Item17:      packet.VoltagePSI(r(ms, 1.0/4, 0.2, 0, 5)),
					Item18:      packet.VoltagePSI(r(ms, 2.0/4, 0.2, 0, 5)),
					Item19:      packet.VoltagePSI(r(ms, 3.0/4, 0.2, 0, 5)),
				},
				State: rdmState,
				// TODO generate data for all fields
				// TODO more realistic data
			})
			ms += 10
		default:
		}
	}
}
