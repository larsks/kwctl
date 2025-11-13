package helpers

import (
	"testing"
)

func TestRangeIterator(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []int
		wantErr  bool
	}{
		{
			name:     "single number",
			input:    "5",
			expected: []int{5},
			wantErr:  false,
		},
		{
			name:     "multiple individual numbers",
			input:    "1,4,10",
			expected: []int{1, 4, 10},
			wantErr:  false,
		},
		{
			name:     "simple range",
			input:    "5-8",
			expected: []int{5, 6, 7, 8},
			wantErr:  false,
		},
		{
			name:     "mixed numbers and ranges",
			input:    "1,4-10",
			expected: []int{1, 4, 5, 6, 7, 8, 9, 10},
			wantErr:  false,
		},
		{
			name:     "multiple ranges",
			input:    "1-3,7-9",
			expected: []int{1, 2, 3, 7, 8, 9},
			wantErr:  false,
		},
		{
			name:     "whitespace around numbers",
			input:    " 1 , 4 , 10 ",
			expected: []int{1, 4, 10},
			wantErr:  false,
		},
		{
			name:     "whitespace around ranges",
			input:    " 5 - 8 ",
			expected: []int{5, 6, 7, 8},
			wantErr:  false,
		},
		{
			name:     "mixed with whitespace",
			input:    " 1 , 4 - 10 ",
			expected: []int{1, 4, 5, 6, 7, 8, 9, 10},
			wantErr:  false,
		},
		{
			name:     "empty string",
			input:    "",
			expected: []int{},
			wantErr:  false,
		},
		{
			name:     "empty parts",
			input:    "1,,4",
			expected: []int{1, 4},
			wantErr:  false,
		},
		{
			name:     "range with same start and end",
			input:    "5-5",
			expected: []int{5},
			wantErr:  false,
		},
		{
			name:     "negative numbers",
			input:    "-5,-2,3",
			expected: []int{-5, -2, 3},
			wantErr:  false,
		},
		{
			name:     "range with negative start",
			input:    "-5-3",
			expected: []int{-5, -4, -3, -2, -1, 0, 1, 2, 3},
			wantErr:  false,
		},
		{
			name:     "range with negative end",
			input:    "5--2",
			expected: []int{},
			wantErr:  true, // descending range
		},
		{
			name:     "range of negative numbers",
			input:    "-10--5",
			expected: []int{-10, -9, -8, -7, -6, -5},
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var results []int
			var errs []error

			for num, err := range RangeIterator(tt.input) {
				if err != nil {
					errs = append(errs, err)
				} else {
					results = append(results, num)
				}
			}

			// Check if we got errors when expected
			if tt.wantErr && len(errs) == 0 {
				t.Errorf("expected error but got none")
			}
			if !tt.wantErr && len(errs) > 0 {
				t.Errorf("unexpected errors: %v", errs)
			}

			// Compare results if we weren't expecting errors
			if !tt.wantErr {
				if len(results) != len(tt.expected) {
					t.Errorf("got %d results, want %d: got %v, want %v",
						len(results), len(tt.expected), results, tt.expected)
					return
				}

				for i := range results {
					if results[i] != tt.expected[i] {
						t.Errorf("at index %d: got %d, want %d (full: got %v, want %v)",
							i, results[i], tt.expected[i], results, tt.expected)
					}
				}
			}
		})
	}
}

func TestRangeIteratorErrors(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		expectErr string
	}{
		{
			name:      "non-numeric value",
			input:     "1,abc,4",
			expectErr: "invalid number",
		},
		{
			name:      "descending range",
			input:     "10-5",
			expectErr: "descending range not allowed",
		},
		{
			name:      "invalid range start",
			input:     "abc-10",
			expectErr: "invalid range start",
		},
		{
			name:      "invalid range end",
			input:     "5-xyz",
			expectErr: "invalid range end",
		},
		{
			name:      "mixed valid and invalid",
			input:     "1,abc,4-6",
			expectErr: "invalid number",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotErr error
			for _, err := range RangeIterator(tt.input) {
				if err != nil {
					gotErr = err
					break
				}
			}

			if gotErr == nil {
				t.Errorf("expected error containing %q but got no error", tt.expectErr)
				return
			}

			if !contains(gotErr.Error(), tt.expectErr) {
				t.Errorf("expected error containing %q, got %q", tt.expectErr, gotErr.Error())
			}
		})
	}
}

func TestRangeIteratorEarlyExit(t *testing.T) {
	input := "1-100"
	count := 0

	for num, err := range RangeIterator(input) {
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		count++
		if num == 5 {
			// Early exit from iteration
			break
		}
	}

	if count != 5 {
		t.Errorf("expected to iterate 5 times, got %d", count)
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && containsHelper(s, substr)))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
