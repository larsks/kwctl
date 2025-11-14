package types

import (
	"fmt"
	"strconv"
	"strings"
)

/*
1 	Memory channel number 3 digit
2 	RX frequency in Hz 10 digit. C clears the channel
3 	RX step size
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
14 	TX frequency in Hz 10 digit, or transmit freq for odd split
15 	TX step size
16 	Lock out
*/

type (
	Channel struct {
		Name      string
		Number    int
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
		TxFreq    int
		TxStep    int
		Lockout   int
	}
)

var EmptyChannel = Channel{}

// ME 101,0145090000,0,0,0,0,0,0,08,08,000,00000000,0,0000000000,0,1
func ParseChannel(s string) (Channel, error) {
	parts := []int{}
	for sval := range strings.SplitSeq(s, ",") {
		ival, err := strconv.Atoi(sval)
		if err != nil {
			return EmptyChannel, err
		}
		parts = append(parts, ival)
	}

	if len(parts) != 16 {
		return EmptyChannel, fmt.Errorf("invalid channel specification")
	}
	return Channel{
		Number:    parts[0],
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
		TxFreq:    parts[13],
		TxStep:    parts[14],
		Lockout:   parts[15],
	}, nil
}

// ME 101,0145090000,0,0,0,0,0,0,08,08,000,00000000,0,0000000000,0,1
func (c Channel) String() string {
	return fmt.Sprintf("[%-6s] ", c.Name) + strings.Join([]string{
		fmt.Sprintf("%03d", c.Number),
		NewFrequencyMHz(&c.RxFreq).String(),
		NewStepSize(&c.RxStep).String(),
		NewShift(&c.Shift).String(),
		NewBool(&c.Reverse).String(),
		NewBool(&c.Tone).String(),
		NewBool(&c.CTCSS).String(),
		NewBool(&c.DCS).String(),
		NewTone(&c.ToneFreq).String(),
		NewTone(&c.CTCSSFreq).String(),
		NewDCS(&c.DCSCode).String(),
		NewFrequencyMHz(&c.Offset).String(),
		NewMode(&c.Mode).String(),
		NewFrequencyMHz(&c.TxFreq).String(),
		NewStepSize(&c.TxStep).String(),
		NewBool(&c.Lockout).String(),
	}, ",")
}

func (c Channel) Serialize() string {
	return fmt.Sprintf("%03d,%010d,%d,%d,%d,%d,%d,%d,%02d,%02d,%03d,%08d,%d,%010d,%d,%d",
		c.Number,
		c.RxFreq,
		c.RxStep,
		c.Shift,
		c.Reverse,
		c.Tone,
		c.CTCSS,
		c.DCS,
		c.ToneFreq,
		c.CTCSSFreq,
		c.DCSCode,
		c.Offset,
		c.Mode,
		c.TxFreq,
		c.TxStep,
		c.Lockout,
	)
}

// RadioSettable interface implementation

func (c *Channel) GetRxFreq() int    { return c.RxFreq }
func (c *Channel) SetRxFreq(v int)   { c.RxFreq = v }
func (c *Channel) GetRxStep() int    { return c.RxStep }
func (c *Channel) SetRxStep(v int)   { c.RxStep = v }
func (c *Channel) GetMode() int      { return c.Mode }
func (c *Channel) SetMode(v int)     { c.Mode = v }
func (c *Channel) GetShift() int     { return c.Shift }
func (c *Channel) SetShift(v int)    { c.Shift = v }
func (c *Channel) GetReverse() int   { return c.Reverse }
func (c *Channel) SetReverse(v int)  { c.Reverse = v }
func (c *Channel) GetOffset() int    { return c.Offset }
func (c *Channel) SetOffset(v int)   { c.Offset = v }
func (c *Channel) GetTone() int      { return c.Tone }
func (c *Channel) SetTone(v int)     { c.Tone = v }
func (c *Channel) GetCTCSS() int     { return c.CTCSS }
func (c *Channel) SetCTCSS(v int)    { c.CTCSS = v }
func (c *Channel) GetDCS() int       { return c.DCS }
func (c *Channel) SetDCS(v int)      { c.DCS = v }
func (c *Channel) GetToneFreq() int  { return c.ToneFreq }
func (c *Channel) SetToneFreq(v int) { c.ToneFreq = v }
func (c *Channel) GetCTCSSFreq() int { return c.CTCSSFreq }
func (c *Channel) SetCTCSSFreq(v int) {
	c.CTCSSFreq = v
}
func (c *Channel) GetDCSCode() int  { return c.DCSCode }
func (c *Channel) SetDCSCode(v int) { c.DCSCode = v }
