package packet

import (
	"github.com/ul-gaul/go-basestation/data/packet/state"
	"time"
)

// Vector represents a 3D vector or coordinates on a 3D plan.
type Vector struct {
	X float32 `csv:"x"`
	Y float32 `csv:"y"`
	Z float32 `csv:"z"`
}

type VoltagePSI float32

// PSI transforms the voltage to PSI. // [0;5] => [0;1000]
func (v VoltagePSI) PSI() float32 {
	return float32(v) * 200
}

type Milliseconds uint64

func (ms Milliseconds) Duration() time.Duration {
	return time.Duration(ms) * time.Millisecond
}

// RocketPacket represents a packet of data received from the rocket.
//
// NOTE: The fields/properties of the struct must be in the same order as the binary data received.
type RocketPacket struct {

	// Time is the number of milliseconds since start
	Time Milliseconds `csv:"time"`

	State state.State `csv:"state"`

	// Temperature is in Celsius (°C). // [-30; 45]
	Temperature float32 `csv:"temperature"`

	// Humidity is a percentage. // [0; 100]
	Humidity float32 `csv:"humidity"`

	// AtmosPressure is in Pascal (Pa). // [0-101325]
	AtmosPressure int32 `csv:"atmos_pressure"`

	// Altitude is in meters (m)
	Altitude float32 `csv:"altitude"`

	// Latitude is the GPS latitude in degree. // [-90;90]
	Latitude float32 `csv:"latitude"`

	// Longitude is the GPS longitude in degree. // [-180;180]
	Longitude float32 `csv:"longitude"`

	// Satellites is the number of GPS satellites the rocket is connected to.
	Satellites uint8 `csv:"nb_sat"`

	// Acceleration is the gravitationnal force experienced by the rocket in g (or m/s²).
	Acceleration Vector `csv:"accel_,inline"`

	// AngularVelocity is the angular speed of the rocket in radian per seconds.
	AngularVelocity Vector `csv:"angular_speed_,inline"`

	// Orientation is the pitch, roll and yaw of the rocket
	Orientation Vector `csv:"orient_,inline"`

	// Magnetism is in millitesla (mT)
	Magnetism Vector `csv:"magnet_,inline"`

	Voltages `csv:"v_,inline"`

	TestBench `csv:"test_,inline"`
}

type TestBench struct {
	LoadCell float32 `csv:"load"`

	// Temperature1 in Celsius (°C) // [0;600]
	Temperature1 float32 `csv:"temp1"`

	// Temperature2 in Celsius (°C) // [0;600]
	Temperature2 float32 `csv:"temp2"`

	// Temperature3 in Celsius (°C) // [0;600]
	Temperature3 float32 `csv:"temp3"`
}

type Voltages struct {
	Rocket      float32    `csv:"rocket"`       // [0;15]
	Tower       float32    `csv:"tower"`        // [0;15]
	IgnitionBox float32    `csv:"ignition_box"` // [0;15]
	Item3       VoltagePSI `csv:"item3"`        // [0;5]
	Item17      VoltagePSI `csv:"item17"`       // [0;5]
	Item18      VoltagePSI `csv:"item18"`       // [0;5]
	Item19      VoltagePSI `csv:"item19"`       // [0;5]
}
