package types

import (
	"testing"
)

func TestFrequencyMHz_Set(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
		wantErr  bool
	}{
		{
			name:     "standard VHF frequency",
			input:    "144.39",
			expected: 144390000,
			wantErr:  false,
		},
		{
			name:     "UHF frequency",
			input:    "446.5",
			expected: 446500000,
			wantErr:  false,
		},
		{
			name:     "whole number frequency",
			input:    "145",
			expected: 145000000,
			wantErr:  false,
		},
		{
			name:     "high precision frequency",
			input:    "145.090",
			expected: 145090000,
			wantErr:  false,
		},
		{
			name:     "zero frequency",
			input:    "0",
			expected: 0,
			wantErr:  false,
		},
		{
			name:     "very small frequency",
			input:    "0.001",
			expected: 1000,
			wantErr:  false,
		},
		{
			name:     "invalid input",
			input:    "not-a-number",
			expected: 0,
			wantErr:  true,
		},
		{
			name:     "negative frequency",
			input:    "-144.39",
			expected: 0,
			wantErr:  true,
		},
		{
			name:     "empty string",
			input:    "",
			expected: 0,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var hz int
			f := NewFrequencyMHz(&hz)
			err := f.Set(tt.input)

			if (err != nil) != tt.wantErr {
				t.Errorf("Set() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && hz != tt.expected {
				t.Errorf("Set() hz = %d, expected %d", hz, tt.expected)
			}
		})
	}
}

func TestFrequencyMHz_String(t *testing.T) {
	tests := []struct {
		name     string
		hz       int
		expected string
	}{
		{
			name:     "standard VHF frequency",
			hz:       144390000,
			expected: "144.390000",
		},
		{
			name:     "UHF frequency",
			hz:       446500000,
			expected: "446.500000",
		},
		{
			name:     "whole MHz value",
			hz:       145000000,
			expected: "145.000000",
		},
		{
			name:     "zero frequency",
			hz:       0,
			expected: "0.000000",
		},
		{
			name:     "high precision",
			hz:       145090000,
			expected: "145.090000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := NewFrequencyMHz(&tt.hz)
			result := f.String()

			if result != tt.expected {
				t.Errorf("String() = %s, expected %s", result, tt.expected)
			}
		})
	}
}

func TestFrequencyMHz_Type(t *testing.T) {
	var hz int
	f := NewFrequencyMHz(&hz)
	expected := "frequencyMHz"

	if f.Type() != expected {
		t.Errorf("Type() = %s, expected %s", f.Type(), expected)
	}
}

func TestFrequencyMHz_RoundTrip(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"VHF", "144.39"},
		{"UHF", "446.5"},
		{"Whole number", "145.0"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var hz int
			f := NewFrequencyMHz(&hz)

			// Set the frequency
			err := f.Set(tt.input)
			if err != nil {
				t.Fatalf("Set() failed: %v", err)
			}

			// Convert back to MHz and set again
			mhzStr := f.String()
			var hz2 int
			f2 := NewFrequencyMHz(&hz2)
			err = f2.Set(mhzStr)
			if err != nil {
				t.Fatalf("Round-trip Set() failed: %v", err)
			}

			// Values should be identical
			if hz != hz2 {
				t.Errorf("Round-trip failed: original %d, after round-trip %d", hz, hz2)
			}
		})
	}
}
