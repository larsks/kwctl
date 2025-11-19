package types

import (
	"fmt"

	"github.com/larsks/kwctl/internal/helpers"
)

type Bool struct {
	valuePtr *int // Pointer to the bool value being configured
}

var boolForward map[int]string = map[int]string{
	0: "false",
	1: "true",
}

var boolReverse map[string]int = helpers.ReverseMap(boolForward)

func NewBool(boolPtr *int) *Bool {
	return &Bool{valuePtr: boolPtr}
}

func (b *Bool) String() string {
	if b.valuePtr == nil {
		return ""
	}
	bool, exists := boolForward[*b.valuePtr]
	if !exists {
		return ""
	}
	return bool
}

// Set parses a bool name and stores it as an integer code
func (b *Bool) Set(value string) error {
	val, exists := boolReverse[value]
	if !exists {
		return fmt.Errorf("invalid bool: %s", value)
	}
	*b.valuePtr = val
	return nil
}

// Type returns the type name for help text
func (b *Bool) Type() string {
	return "bool"
}
