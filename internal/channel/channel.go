package channel

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
		Channel   int
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
		TxFreq    int
		TxStep    int
		Lock      int
	}
)

// ME 101,0145090000,0,0,0,0,0,0,08,08,000,00000000,0,0000000000,0,1
func ParseChannel(s string) (Channel, error) {
	parts := []int{}
	for sval := range strings.SplitSeq(s, ",") {
		ival, err := strconv.Atoi(sval)
		if err != nil {
			return Channel{}, err
		}
		parts = append(parts, ival)
	}

	if len(parts) != 16 {
		return Channel{}, fmt.Errorf("invalid channel specification")
	}
	return Channel{
		Channel:   parts[0],
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
		TxFreq:    parts[13],
		TxStep:    parts[14],
		Lock:      parts[15],
	}, nil
}

// ME 101,0145090000,0,0,0,0,0,0,08,08,000,00000000,0,0000000000,0,1
func (c Channel) String() string {
	return fmt.Sprintf("%03d,%010d,%d,%d,%d,%d,%d,%d,%02d,%02d,%03d,%08d,%d,%010d,%d,%d",
		c.Channel,
		c.RxFreq,
		c.RxStep,
		c.Shift,
		c.Reverse,
		c.Tone,
		c.CTCSS,
		c.DCS,
		c.ToneFreq,
		c.CTCSSFreq,
		c.DCSFreq,
		c.Offset,
		c.Mode,
		c.TxFreq,
		c.TxStep,
		c.Lock,
	)
}
