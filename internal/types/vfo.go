package types

import (
	"encoding/json"
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
		DCSFreq   int
		Offset    int
		Mode      int
	}
)

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
		DCSFreq:   parts[10],
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
		v.DCSFreq,
		v.Offset,
		v.Mode,
	)
}

// MarshalJSON implements custom JSON marshaling with human-friendly values
func (v VFO) MarshalJSON() ([]byte, error) {
	// Create a struct with human-friendly representations
	humanFriendly := struct {
		VFO       int    `json:"VFO"`
		RxFreq    string `json:"RxFreq"`
		RxStep    string `json:"RxStep"`
		Shift     int    `json:"Shift"`
		Reverse   int    `json:"Reverse"`
		Tone      int    `json:"Tone"`
		CTCSS     int    `json:"CTCSS"`
		DCS       int    `json:"DCS"`
		ToneFreq  int    `json:"ToneFreq"`
		CTCSSFreq int    `json:"CTCSSFreq"`
		DCSFreq   int    `json:"DCSFreq"`
		Offset    string `json:"Offset"`
		Mode      string `json:"Mode"`
	}{
		VFO:       v.VFO,
		RxFreq:    NewFrequencyMHz(&v.RxFreq).String(),
		RxStep:    NewStepSize(&v.RxStep).String(),
		Shift:     v.Shift,
		Reverse:   v.Reverse,
		Tone:      v.Tone,
		CTCSS:     v.CTCSS,
		DCS:       v.DCS,
		ToneFreq:  v.ToneFreq,
		CTCSSFreq: v.CTCSSFreq,
		DCSFreq:   v.DCSFreq,
		Offset:    NewFrequencyMHz(&v.Offset).String(),
		Mode:      NewMode(&v.Mode).String(),
	}

	return json.Marshal(humanFriendly)
}

// FO 1,0145090000,0,0,0,0,0,0,08,08,000,00000000,0
func (v VFO) String() string {
	s, _ := json.MarshalIndent(v, "", "  ")
	return string(s)
}
