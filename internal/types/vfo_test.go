package types

import (
	"encoding/json"
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

func TestVFO_MarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		vfo      VFO
		expected map[string]any
	}{
		{
			name: "standard VFO configuration",
			vfo:  VFO{1, 145090000, 0, 0, 0, 0, 0, 0, 8, 8, 0, 0, 0},
			expected: map[string]any{
				"VFO":       float64(1),
				"RxFreq":    "145.090000",
				"RxStep":    "5",
				"Shift":     float64(0),
				"Reverse":   float64(0),
				"Tone":      float64(0),
				"CTCSS":     float64(0),
				"DCS":       float64(0),
				"ToneFreq":  float64(8),
				"CTCSSFreq": float64(8),
				"DCSFreq":   float64(0),
				"Offset":    "0.000000",
				"Mode":      "FM",
			},
		},
		{
			name: "UHF with offset and NFM mode",
			vfo:  VFO{0, 446500000, 4, 1, 0, 0, 0, 0, 8, 8, 0, 5000000, 1},
			expected: map[string]any{
				"VFO":       float64(0),
				"RxFreq":    "446.500000",
				"RxStep":    "12.5",
				"Shift":     float64(1),
				"Reverse":   float64(0),
				"Tone":      float64(0),
				"CTCSS":     float64(0),
				"DCS":       float64(0),
				"ToneFreq":  float64(8),
				"CTCSSFreq": float64(8),
				"DCSFreq":   float64(0),
				"Offset":    "5.000000",
				"Mode":      "NFM",
			},
		},
		{
			name: "AM mode with different step size",
			vfo:  VFO{1, 118000000, 7, 0, 0, 0, 0, 0, 8, 8, 0, 0, 2},
			expected: map[string]any{
				"VFO":       float64(1),
				"RxFreq":    "118.000000",
				"RxStep":    "25",
				"Shift":     float64(0),
				"Reverse":   float64(0),
				"Tone":      float64(0),
				"CTCSS":     float64(0),
				"DCS":       float64(0),
				"ToneFreq":  float64(8),
				"CTCSSFreq": float64(8),
				"DCSFreq":   float64(0),
				"Offset":    "0.000000",
				"Mode":      "AM",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonBytes, err := json.Marshal(tt.vfo)
			if err != nil {
				t.Fatalf("MarshalJSON() failed: %v", err)
			}

			var result map[string]any
			err = json.Unmarshal(jsonBytes, &result)
			if err != nil {
				t.Fatalf("Unmarshal failed: %v", err)
			}

			// Check each expected field
			for key, expectedVal := range tt.expected {
				actualVal, exists := result[key]
				if !exists {
					t.Errorf("Missing field %s in JSON output", key)
					continue
				}

				// Compare values (handles both string and numeric types)
				if actualVal != expectedVal {
					t.Errorf("Field %s: got %v, expected %v", key, actualVal, expectedVal)
				}
			}
		})
	}
}

func TestVFO_String(t *testing.T) {
	vfo := VFO{1, 145090000, 0, 0, 0, 0, 0, 0, 8, 8, 0, 0, 0}
	result := vfo.String()

	// Should produce valid JSON
	var parsed map[string]any
	err := json.Unmarshal([]byte(result), &parsed)
	if err != nil {
		t.Fatalf("String() did not produce valid JSON: %v", err)
	}

	// Should contain human-friendly values
	if parsed["RxFreq"] != "145.090000" {
		t.Errorf("RxFreq not human-friendly: got %v", parsed["RxFreq"])
	}
	if parsed["Mode"] != "FM" {
		t.Errorf("Mode not human-friendly: got %v", parsed["Mode"])
	}
	if parsed["RxStep"] != "5" {
		t.Errorf("RxStep not human-friendly: got %v", parsed["RxStep"])
	}
}
