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
			Channel{101, 145090000, 0, 0, 0, 0, 0, 0, 8, 8, 0, 0, 0, 0, 0, 1},
			true,
		},
		{
			"001,0146820000,0,2,0,1,0,0,23,23,000,00600000,0,0000000000,0,0",
			Channel{1, 146820000, 0, 2, 0, 1, 0, 0, 23, 23, 0, 600000, 0, 0, 0, 0},
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
			"101,0145090000,0,0,0,0,0,0,08,08,000,00000000,0,0000000000,0,1",
			Channel{101, 145090000, 0, 0, 0, 0, 0, 0, 8, 8, 0, 0, 0, 0, 0, 1},
			true,
		},
		{
			"001,0145090000,0,0,0,0,0,0,08,08,000,00000000,0,0000000000,0,1",
			Channel{1, 145090000, 0, 0, 0, 0, 0, 0, 8, 8, 0, 0, 0, 0, 0, 1},
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
