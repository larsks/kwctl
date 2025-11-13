package types

import (
	"fmt"
	"strconv"
)

// FrequencyMHz is a pflag.Value implementation that handles frequency conversion
// from MHz (human-readable format) to Hz (internal storage format).
type FrequencyMHz struct {
	valuePtr *int // Pointer to the Hz value being configured
}

// NewFrequencyMHz creates a new FrequencyMHz flag value that updates the provided Hz pointer
func NewFrequencyMHz(hzPtr *int) *FrequencyMHz {
	return &FrequencyMHz{valuePtr: hzPtr}
}

// String returns the current frequency value in MHz format
func (f *FrequencyMHz) String() string {
	if f.valuePtr == nil {
		return "0"
	}
	mhz := float64(*f.valuePtr) / 1_000_000
	return fmt.Sprintf("%.6f", mhz)
}

// Set parses a MHz value and stores it as Hz
func (f *FrequencyMHz) Set(value string) error {
	mhz, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return fmt.Errorf("invalid frequency: %w", err)
	}
	if mhz < 0 {
		return fmt.Errorf("frequency cannot be negative")
	}
	hz := int(mhz * 1_000_000)
	*f.valuePtr = hz
	return nil
}

// Type returns the type name for help text
func (f *FrequencyMHz) Type() string {
	return "frequencyMHz"
}
