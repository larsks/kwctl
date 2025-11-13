package types

import (
	"fmt"

	"github.com/larsks/kwctl/internal/helpers"
)

type Mode struct {
	valuePtr *int // Pointer to the mode value being configured
}

var modeForward map[int]string = map[int]string{
	0: "FM",
	1: "NFM",
	2: "AM",
}

var modeReverse map[string]int = helpers.ReverseMap(modeForward)

func NewMode(modePtr *int) *Mode {
	return &Mode{valuePtr: modePtr}
}

func (m *Mode) String() string {
	if m.valuePtr == nil {
		return ""
	}
	mode, exists := modeForward[*m.valuePtr]
	if !exists {
		return ""
	}
	return mode
}

// Set parses a mode name and stores it as an integer code
func (m *Mode) Set(value string) error {
	val, exists := modeReverse[value]
	if !exists {
		return fmt.Errorf("invalid mode: %s", value)
	}
	*m.valuePtr = val
	return nil
}

// Type returns the type name for help text
func (m *Mode) Type() string {
	return "mode"
}
