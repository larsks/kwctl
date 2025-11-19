package types

import (
	"fmt"
	"strconv"
	"strings"
)

/*
1 	Band
2 	Frequency in Hz 10 digit. must be within selected band
3 	Step size
4 	Shift direction
5 	Reverse
6 	Tone status
7 	CTCSS status
8 	DCS status
9 	Tone frequency
10 	CTCSS frequency
11 	DCS frequency
12 	Offset frequency in Hz 8 digit
13 	Mode
*/

type (
	DisplayVFO struct {
		VFO       int
		RxFreq    string
		RxStep    string
		Shift     string
		Reverse   string
		Tone      string
		CTCSS     string
		DCS       string
		ToneFreq  string
		CTCSSFreq string
		DCSCode   string
		Offset    string
		Mode      string
	}

	VFO struct {
		VFO       int
		RxFreq    int
		RxStep    int
		Shift     int
		Reverse   int
		Tone      int
		CTCSS     int
		DCS       int
		ToneFreq  int
		CTCSSFreq int
		DCSCode   int
		Offset    int
		Mode      int
	}
)

var EmptyVFO = VFO{}

// FO 1,0145090000,0,0,0,0,0,0,08,08,000,00000000,0
func ParseVFO(s string) (VFO, error) {
	parts := []int{}
	for sval := range strings.SplitSeq(s, ",") {
		ival, err := strconv.Atoi(sval)
		if err != nil {
			return VFO{}, err
		}
		parts = append(parts, ival)
	}

	if len(parts) != 13 {
		return VFO{}, fmt.Errorf("invalid vfo specification")
	}
	return VFO{
		VFO:       parts[0],
		RxFreq:    parts[1],
		RxStep:    parts[2],
		Shift:     parts[3],
		Reverse:   parts[4],
		Tone:      parts[5],
		CTCSS:     parts[6],
		DCS:       parts[7],
		ToneFreq:  parts[8],
		CTCSSFreq: parts[9],
		DCSCode:   parts[10],
		Offset:    parts[11],
		Mode:      parts[12],
	}, nil
}

// FO 1,0145090000,0,0,0,0,0,0,08,08,000,00000000,0
func (v VFO) Serialize() string {
	return fmt.Sprintf("%d,%010d,%d,%d,%d,%d,%d,%d,%02d,%02d,%03d,%08d,%d",
		v.VFO,
		v.RxFreq,
		v.RxStep,
		v.Shift,
		v.Reverse,
		v.Tone,
		v.CTCSS,
		v.DCS,
		v.ToneFreq,
		v.CTCSSFreq,
		v.DCSCode,
		v.Offset,
		v.Mode,
	)
}

func (v VFO) Values() []string {
	return []string{
		fmt.Sprintf("%d", v.VFO),
		NewFrequencyMHz(&v.RxFreq).String(),
		NewStepSize(&v.RxStep).String(),
		NewShift(&v.Shift).String(),
		NewBool(&v.Reverse).String(),
		NewBool(&v.Tone).String(),
		NewBool(&v.CTCSS).String(),
		NewBool(&v.DCS).String(),
		NewTone(&v.ToneFreq).String(),
		NewTone(&v.CTCSSFreq).String(),
		NewDCS(&v.DCSCode).String(),
		NewFrequencyMHz(&v.Offset).String(),
		NewMode(&v.Mode).String(),
	}
}

func (v VFO) Display() DisplayVFO {
	return DisplayVFO{
		VFO:       v.VFO,
		RxFreq:    NewFrequencyMHz(&v.RxFreq).String(),
		RxStep:    NewStepSize(&v.RxStep).String(),
		Shift:     NewShift(&v.Shift).String(),
		Reverse:   NewBool(&v.Reverse).String(),
		Tone:      NewBool(&v.Tone).String(),
		CTCSS:     NewBool(&v.CTCSS).String(),
		DCS:       NewBool(&v.DCS).String(),
		ToneFreq:  NewTone(&v.ToneFreq).String(),
		CTCSSFreq: NewTone(&v.CTCSSFreq).String(),
		DCSCode:   NewDCS(&v.DCSCode).String(),
		Offset:    NewFrequencyMHz(&v.Offset).String(),
		Mode:      NewMode(&v.Mode).String(),
	}
}

func (v VFO) String() string {
	return strings.Join([]string{
		fmt.Sprintf("%d", v.VFO),
		NewFrequencyMHz(&v.RxFreq).String(),
		NewStepSize(&v.RxStep).String(),
		NewShift(&v.Shift).String(),
		NewBool(&v.Reverse).String(),
		NewBool(&v.Tone).String(),
		NewBool(&v.CTCSS).String(),
		NewBool(&v.DCS).String(),
		NewTone(&v.ToneFreq).String(),
		NewTone(&v.CTCSSFreq).String(),
		NewDCS(&v.DCSCode).String(),
		NewFrequencyMHz(&v.Offset).String(),
		NewMode(&v.Mode).String(),
	}, ",")
}

// RadioSettable interface implementation

func (v *VFO) GetRxFreq() int     { return v.RxFreq }
func (v *VFO) SetRxFreq(val int)  { v.RxFreq = val }
func (v *VFO) GetRxStep() int     { return v.RxStep }
func (v *VFO) SetRxStep(val int)  { v.RxStep = val }
func (v *VFO) GetMode() int       { return v.Mode }
func (v *VFO) SetMode(val int)    { v.Mode = val }
func (v *VFO) GetShift() int      { return v.Shift }
func (v *VFO) SetShift(val int)   { v.Shift = val }
func (v *VFO) GetReverse() int    { return v.Reverse }
func (v *VFO) SetReverse(val int) { v.Reverse = val }
func (v *VFO) GetOffset() int     { return v.Offset }
func (v *VFO) SetOffset(val int)  { v.Offset = val }
func (v *VFO) GetTone() int       { return v.Tone }
func (v *VFO) SetTone(val int)    { v.Tone = val }
func (v *VFO) GetCTCSS() int      { return v.CTCSS }
func (v *VFO) SetCTCSS(val int)   { v.CTCSS = val }
func (v *VFO) GetDCS() int        { return v.DCS }
func (v *VFO) SetDCS(val int)     { v.DCS = val }
func (v *VFO) GetToneFreq() int   { return v.ToneFreq }
func (v *VFO) SetToneFreq(val int) {
	v.ToneFreq = val
}
func (v *VFO) GetCTCSSFreq() int { return v.CTCSSFreq }
func (v *VFO) SetCTCSSFreq(val int) {
	v.CTCSSFreq = val
}
func (v *VFO) GetDCSCode() int    { return v.DCSCode }
func (v *VFO) SetDCSCode(val int) { v.DCSCode = val }
