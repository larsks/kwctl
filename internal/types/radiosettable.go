package types

// RadioSettable defines an interface for types that can have radio settings
// (frequency, mode, shift, tones, etc.) configured via common flags.
// Both Channel and VFO implement this interface.
type RadioSettable interface {
	GetRxFreq() int
	SetRxFreq(int)

	GetRxStep() int
	SetRxStep(int)

	GetMode() int
	SetMode(int)

	GetShift() int
	SetShift(int)

	GetReverse() int
	SetReverse(int)

	GetOffset() int
	SetOffset(int)

	GetTone() int
	SetTone(int)

	GetCTCSS() int
	SetCTCSS(int)

	GetDCS() int
	SetDCS(int)

	GetToneFreq() int
	SetToneFreq(int)

	GetCTCSSFreq() int
	SetCTCSSFreq(int)

	GetDCSCode() int
	SetDCSCode(int)
}
