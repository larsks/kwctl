package types

/*
00 	67.0
01 	69.3
02 	71.9
03 	74.4
04 	77.0
05 	79.7
06 	82.5
07 	85.4
08 	88.5
19 	91.5
10 	94.8
11 	97.4
12 	100.0
13 	103.5
14 	107.2
15 	110.9
16 	114.8
17 	118.8
18 	123.0
19 	127.3
20 	131.8
21 	136.5
22 	141.3
23 	146.2
24 	151.4
25 	156.7
26 	162.2
27 	167.9
28 	173.8
29 	179.9
30 	186.2
31 	192.8
32 	203.5
33 	206.5
34 	210.7
35 	218.1
36 	225.7
37 	229.1
38 	233.6
39 	241.8
40 	250.3
41 	254.1
*/

import (
	"fmt"

	"github.com/larsks/kwctl/internal/helpers"
)

type Tone struct {
	valuePtr *int // Pointer to the tone value being configured
}

var toneForward map[int]string = map[int]string{
	0:  "67.0",
	1:  "69.3",
	2:  "71.9",
	3:  "74.4",
	4:  "77.0",
	5:  "79.7",
	6:  "82.5",
	7:  "85.4",
	8:  "88.5",
	9:  "91.5",
	10: "94.8",
	11: "97.4",
	12: "100.0",
	13: "103.5",
	14: "107.2",
	15: "110.9",
	16: "114.8",
	17: "118.8",
	18: "123.0",
	19: "127.3",
	20: "131.8",
	21: "136.5",
	22: "141.3",
	23: "146.2",
	24: "151.4",
	25: "156.7",
	26: "162.2",
	27: "167.9",
	28: "173.8",
	29: "179.9",
	30: "186.2",
	31: "192.8",
	32: "203.5",
	33: "206.5",
	34: "210.7",
	35: "218.1",
	36: "225.7",
	37: "229.1",
	38: "233.6",
	39: "241.8",
	40: "250.3",
	41: "254.1",
}

var toneReverse map[string]int = helpers.ReverseMap(toneForward)

func NewTone(tonePtr *int) *Tone {
	return &Tone{valuePtr: tonePtr}
}

func (t *Tone) String() string {
	if t.valuePtr == nil {
		return ""
	}
	tone, exists := toneForward[*t.valuePtr]
	if !exists {
		return ""
	}
	return tone
}

// Set parses a tone name and stores it as an integer code
func (t *Tone) Set(value string) error {
	val, exists := toneReverse[value]
	if !exists {
		return fmt.Errorf("invalid tone: %s", value)
	}
	*t.valuePtr = val
	return nil
}

// Type returns the type name for help text
func (t *Tone) Type() string {
	return "tone"
}
