package types

import (
	"testing"
)

type (
	ChannelTestItem struct {
		raw    string
		parsed Channel
		valid  bool
	}
)

func TestParseChannel(t *testing.T) {
	inputs := []ChannelTestItem{
		{
			"101,0145090000,0,0,0,0,0,0,08,08,000,00000000,0,0000000000,0,1",
			Channel{Number: 101, RxFreq: 145090000, RxStep: 0, Shift: 0, Reverse: 0, Tone: 0, CTCSS: 0, DCS: 0, ToneFreq: 8, CTCSSFreq: 8, DCSCode: 0, Offset: 0, Mode: 0, TxFreq: 0, TxStep: 0, Lockout: 1},
			true,
		},
		{
			"001,0146820000,0,2,0,1,0,0,23,23,000,00600000,0,0000000000,0,0",
			Channel{Number: 1, RxFreq: 146820000, RxStep: 0, Shift: 2, Reverse: 0, Tone: 1, CTCSS: 0, DCS: 0, ToneFreq: 23, CTCSSFreq: 23, DCSCode: 0, Offset: 600000, Mode: 0, TxFreq: 0, TxStep: 0, Lockout: 0},
			true,
		},
		{
			"0",
			Channel{},
			false,
		},
	}

	for _, input := range inputs {
		have, err := ParseChannel(input.raw)
		if err != nil && input.valid {
			t.Errorf("expected success, failed with: %v", err)
		} else if err == nil && !input.valid {
			t.Errorf("expected error")
		}

		if input.valid && have != input.parsed {
			t.Errorf("unexpected result")
		}
	}
}

func TestStringifyChannel(t *testing.T) {
	inputs := []ChannelTestItem{
		{
			"[      ] 101,145.09,0,simplex,false,false,false,false,0.0,0.0,023,0.000,FM,0.00,0,true",
			Channel{Number: 101, RxFreq: 145090000, RxStep: 0, Shift: 0, Reverse: 0, Tone: 0, CTCSS: 0, DCS: 0, ToneFreq: 8, CTCSSFreq: 8, DCSCode: 0, Offset: 0, Mode: 0, TxFreq: 0, TxStep: 0, Lockout: 1},
			true,
		},
		{
			"[      ] 001,145.09,0,simplex,false,false,false,false,0.0,0.0,023,0.000,FM,0.00,0,false",
			Channel{Number: 1, RxFreq: 145090000, RxStep: 0, Shift: 0, Reverse: 0, Tone: 0, CTCSS: 0, DCS: 0, ToneFreq: 8, CTCSSFreq: 8, DCSCode: 0, Offset: 0, Mode: 0, TxFreq: 0, TxStep: 0, Lockout: 0},
			true,
		},
	}

	for _, input := range inputs {
		raw := input.parsed.String()
		if raw != input.raw {
			t.Errorf("have %s, expected %s", raw, input.raw)
		}
	}
}
