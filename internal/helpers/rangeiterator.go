package helpers

import (
	"fmt"
	"iter"
	"strconv"
	"strings"
)

// RangeIterator parses a comma-separated string of numbers and ranges,
// yielding individual integer values. Supports syntax like "1,4,10" or "1,4-10".
//
// Examples:
//   - Individual numbers: "1,4,10" yields 1, 4, 10
//   - Ranges: "5-8" yields 5, 6, 7, 8
//   - Mixed: "1,4-10" yields 1, 4, 5, 6, 7, 8, 9, 10
//
// Whitespace around numbers and the range operator is automatically trimmed.
// Ranges must be ascending (start <= end); descending ranges yield an error.
// Invalid input (non-numeric, malformed syntax) yields an error for that element.
//
// The iterator yields both a value and an error. When an error occurs for a
// specific element, the value will be 0 and the error will be non-nil.
// The iterator continues processing remaining elements after an error.
func RangeIterator(input string) iter.Seq2[int, error] {
	return func(yield func(int, error) bool) {
		// Split by comma
		parts := strings.Split(input, ",")

		for _, part := range parts {
			part = strings.TrimSpace(part)

			// Skip empty parts
			if part == "" {
				continue
			}

			// Look for "-" that could be a range operator
			// Must not be at position 0 (negative number) or after another "-"
			dashIdx := -1
			for i := 1; i < len(part); i++ {
				if part[i] == '-' && part[i-1] != '-' {
					dashIdx = i
					break
				}
			}

			if dashIdx > 0 {
				// This is a range
				startStr := strings.TrimSpace(part[:dashIdx])
				endStr := strings.TrimSpace(part[dashIdx+1:])

				// Parse start
				start, err := strconv.Atoi(startStr)
				if err != nil {
					if !yield(0, fmt.Errorf("invalid range start %q: %w", startStr, err)) {
						return
					}
					continue
				}

				// Parse end
				end, err := strconv.Atoi(endStr)
				if err != nil {
					if !yield(0, fmt.Errorf("invalid range end %q: %w", endStr, err)) {
						return
					}
					continue
				}

				// Validate ascending
				if start > end {
					if !yield(0, fmt.Errorf("descending range not allowed: %d-%d", start, end)) {
						return
					}
					continue
				}

				// Yield all values in range (inclusive)
				for i := start; i <= end; i++ {
					if !yield(i, nil) {
						return
					}
				}
			} else {
				// Single number
				num, err := strconv.Atoi(part)
				if err != nil {
					if !yield(0, fmt.Errorf("invalid number %q: %w", part, err)) {
						return
					}
					continue
				}

				if !yield(num, nil) {
					return
				}
			}
		}
	}
}
