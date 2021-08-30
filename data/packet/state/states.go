package state

// State represents the state of the rocket.
// Each bit of State represent a status (on/off).
type State uint8

type Item uint8

const (
	I1 Item = 0b0000_0001 << iota
	I2
	I3
	I4
	I5
	I6
	I7
	I8
)

func (s State) All() []bool {
	slice := make([]bool, 8)
	for i := 0; i < 8; i++ {
		slice[i] = (s & (0b0000_0001 << i)) > 0
	}
	return slice
}

func (s State) StatusOf(i Item) bool {
	return (s & State(i)) > 0
}

// I1 returns the status (on/off) represented by the rightmost bit
func (s State) I1() bool {
	return s.StatusOf(I1)
}

// I2 returns the status (on/off) represented by the 2nd rightmost bit
func (s State) I2() bool {
	return s.StatusOf(I2)
}

// I3 returns the status (on/off) represented by the 3rd rightmost bit
func (s State) I3() bool {
	return s.StatusOf(I3)
}

// I4 returns the status (on/off) represented by the 4th rightmost bit
func (s State) I4() bool {
	return s.StatusOf(I4)
}

// I5 returns the status (on/off) represented by the 5th rightmost bit
func (s State) I5() bool {
	return s.StatusOf(I5)
}

// I6 returns the status (on/off) represented by the 6th rightmost bit
func (s State) I6() bool {
	return s.StatusOf(I6)
}

// I7 returns the status (on/off) represented by the 7th rightmost bit
func (s State) I7() bool {
	return s.StatusOf(I7)
}

// I8 returns the status (on/off) represented by the 8th rightmost bit
func (s State) I8() bool {
	return s.StatusOf(I8)
}

// WithStatusOf returns a copy with the specified item set to the specified status
func (s State) WithStatusOf(i Item, status bool) State {
	if status {
		return s | State(i)
	}
	return s & ^State(i)
}

// WithI1 returns a copy with I1 set to the specified status
func (s State) WithI1(status bool) State {
	return s.WithStatusOf(I1, status)
}

// WithI2 returns a copy with I2 set to the specified status.
func (s State) WithI2(status bool) State {
	return s.WithStatusOf(I2, status)
}

// WithI3 returns a copy with I3 set to the specified status.
func (s State) WithI3(status bool) State {
	return s.WithStatusOf(I3, status)
}

// WithI4 returns a copy with I4 set to the specified status.
func (s State) WithI4(status bool) State {
	return s.WithStatusOf(I4, status)
}

// WithI5 returns a copy with I5 set to the specified status.
func (s State) WithI5(status bool) State {
	return s.WithStatusOf(I5, status)
}

// WithI6 returns a copy with I6 set to the specified status.
func (s State) WithI6(status bool) State {
	return s.WithStatusOf(I6, status)
}

// WithI7 returns a copy with I7 set to the specified status.
func (s State) WithI7(status bool) State {
	return s.WithStatusOf(I7, status)
}

// WithI8 returns a copy with I8 set to the specified status.
func (s State) WithI8(status bool) State {
	return s.WithStatusOf(I8, status)
}
