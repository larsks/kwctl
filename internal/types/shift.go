package types

import (
	"fmt"

	"github.com/larsks/kwctl/internal/helpers"
)

type Shift struct {
	valuePtr *int // Pointer to the shift value being configured
}

var shiftForward map[int]string = map[int]string{
	0: "simplex",
	1: "up",
	2: "down",
}

var shiftReverse map[string]int = helpers.ReverseMap(shiftForward)

func NewShift(shiftPtr *int) *Shift {
	return &Shift{valuePtr: shiftPtr}
}

func (s *Shift) String() string {
	if s.valuePtr == nil {
		return ""
	}
	shift, exists := shiftForward[*s.valuePtr]
	if !exists {
		return ""
	}
	return shift
}

// Set parses a shift name and stores it as an integer code
func (m *Shift) Set(value string) error {
	val, exists := shiftReverse[value]
	if !exists {
		return fmt.Errorf("invalid shift: %s", value)
	}
	*m.valuePtr = val
	return nil
}

// Type returns the type name for help text
func (m *Shift) Type() string {
	return "shift"
}
