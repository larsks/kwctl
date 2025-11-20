package types

import (
	"fmt"
)

type (
	BandMode int
)

const (
	BAND_MODE_DUAL   BandMode = 0
	BAND_MODE_SINGLE BandMode = 1
)

var (
	bandModeNames = map[string]BandMode{
		"dual":   BAND_MODE_DUAL,
		"single": BAND_MODE_SINGLE,
	}
)

func (v BandMode) String() string {
	switch v {
	case BAND_MODE_DUAL:
		return "dual"
	case BAND_MODE_SINGLE:
		return "single"
	default:
		return "<invalid>"
	}
}

func ParseBandMode(s string) (BandMode, error) {
	if val, exists := bandModeNames[s]; exists {
		return val, nil
	}

	return 0, fmt.Errorf("invalid band mode: %s", s)
}
