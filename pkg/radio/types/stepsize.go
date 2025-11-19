package types

import (
	"fmt"

	"github.com/larsks/kwctl/internal/helpers"
)

type StepSize struct {
	valuePtr *int // Pointer to the Hz value being configured
}

var stepSizeForward map[int]string = map[int]string{
	0x0: "5",
	0x1: "6.25",
	0x2: "28.33",
	0x3: "10",
	0x4: "12.5",
	0x5: "15",
	0x6: "20",
	0x7: "25",
	0x8: "30",
	0x9: "50",
	0xA: "100",
}

var stepSizeReverse map[string]int = helpers.ReverseMap(stepSizeForward)

func NewStepSize(hzPtr *int) *StepSize {
	return &StepSize{valuePtr: hzPtr}
}

func (f *StepSize) String() string {
	if f.valuePtr == nil {
		return "0"
	}
	hz, exists := stepSizeForward[*f.valuePtr]
	if !exists {
		return "0"
	}
	return hz
}

// Set parses a MHz value and stores it as Hz
func (f *StepSize) Set(value string) error {
	val, exists := stepSizeReverse[value]
	if !exists {
		return fmt.Errorf("invalid step size: %s", value)
	}
	*f.valuePtr = val
	return nil
}

// Type returns the type name for help text
func (f *StepSize) Type() string {
	return "stepSize"
}
