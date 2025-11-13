package types

import (
	"testing"
)

func TestMode_Set(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
		wantErr  bool
	}{
		{
			name:     "FM mode",
			input:    "FM",
			expected: 0,
			wantErr:  false,
		},
		{
			name:     "NFM mode",
			input:    "NFM",
			expected: 1,
			wantErr:  false,
		},
		{
			name:     "AM mode",
			input:    "AM",
			expected: 2,
			wantErr:  false,
		},
		{
			name:     "invalid mode",
			input:    "SSB",
			expected: 0,
			wantErr:  true,
		},
		{
			name:     "lowercase input",
			input:    "fm",
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
			var mode int
			m := NewMode(&mode)
			err := m.Set(tt.input)

			if (err != nil) != tt.wantErr {
				t.Errorf("Set() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && mode != tt.expected {
				t.Errorf("Set() mode = %d, expected %d", mode, tt.expected)
			}
		})
	}
}

func TestMode_String(t *testing.T) {
	tests := []struct {
		name     string
		mode     int
		expected string
	}{
		{
			name:     "FM mode",
			mode:     0,
			expected: "FM",
		},
		{
			name:     "NFM mode",
			mode:     1,
			expected: "NFM",
		},
		{
			name:     "AM mode",
			mode:     2,
			expected: "AM",
		},
		{
			name:     "invalid mode code",
			mode:     99,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMode(&tt.mode)
			result := m.String()

			if result != tt.expected {
				t.Errorf("String() = %s, expected %s", result, tt.expected)
			}
		})
	}
}

func TestMode_Type(t *testing.T) {
	var mode int
	m := NewMode(&mode)
	expected := "mode"

	if m.Type() != expected {
		t.Errorf("Type() = %s, expected %s", m.Type(), expected)
	}
}

func TestMode_RoundTrip(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"FM", "FM"},
		{"NFM", "NFM"},
		{"AM", "AM"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var mode int
			m := NewMode(&mode)

			// Set the mode
			err := m.Set(tt.input)
			if err != nil {
				t.Fatalf("Set() failed: %v", err)
			}

			// Convert back to string and set again
			modeStr := m.String()
			var mode2 int
			m2 := NewMode(&mode2)
			err = m2.Set(modeStr)
			if err != nil {
				t.Fatalf("Round-trip Set() failed: %v", err)
			}

			// Values should be identical
			if mode != mode2 {
				t.Errorf("Round-trip failed: original %d, after round-trip %d", mode, mode2)
			}
		})
	}
}
