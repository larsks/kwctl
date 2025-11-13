package types

import (
	"testing"
)

type (
	VFOTestItem struct {
		raw    string
		parsed VFO
		valid  bool
	}
)

func TestParseVFO(t *testing.T) {
	inputs := []VFOTestItem{
		{
			"1,0145090000,0,0,0,0,0,0,08,08,000,00000000,0",
			VFO{1, 145090000, 0, 0, 0, 0, 0, 0, 8, 8, 0, 0, 0},
			true,
		},
	}

	for _, input := range inputs {
		have, err := ParseVFO(input.raw)
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

func TestSerializeVFO(t *testing.T) {
	inputs := []VFOTestItem{
		{
			"1,0145090000,0,0,0,0,0,0,08,08,000,00000000,0",
			VFO{1, 145090000, 0, 0, 0, 0, 0, 0, 8, 8, 0, 0, 0},
			true,
		},
	}

	for _, input := range inputs {
		raw := input.parsed.Serialize()
		if raw != input.raw {
			t.Errorf("have %s, expected %s", raw, input.raw)
		}
	}
}
