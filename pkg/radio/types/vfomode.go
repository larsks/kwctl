package types

import (
	"fmt"
)

type (
	VfoMode int
)

const (
	VFO_MODE_VFO    VfoMode = 0
	VFO_MODE_MEMORY VfoMode = 1
	VFO_MODE_CALL   VfoMode = 2
	VFO_MODE_WX     VfoMode = 3
)

var (
	vfoModeNames = map[string]VfoMode{
		"vfo":    VFO_MODE_VFO,
		"memory": VFO_MODE_MEMORY,
		"call":   VFO_MODE_CALL,
		"wx":     VFO_MODE_WX,
	}
)

func (v VfoMode) String() string {
	switch v {
	case VFO_MODE_VFO:
		return "vfo"
	case VFO_MODE_MEMORY:
		return "memory"
	case VFO_MODE_CALL:
		return "call"
	case VFO_MODE_WX:
		return "wx"
	default:
		return "<invalid>"
	}
}

func ParseVfoMode(s string) (VfoMode, error) {
	if val, exists := vfoModeNames[s]; exists {
		return val, nil
	}

	return 0, fmt.Errorf("invalid vfo mode: %s", s)
}
