package helpers

import "strings"

// Unindent removes common leading whitespace from a multiline string.
//
// The function performs the following operations:
//  1. Removes any leading blank lines
//  2. Removes any trailing blank lines
//  3. Finds the minimum leading whitespace (spaces and tabs) among all non-empty lines
//  4. Removes that amount of leading whitespace from each line
//  5. Replaces lines containing only whitespace with blank lines
//  6. Preserves trailing newline if the input ended with one
//
// For empty strings or strings containing only whitespace, returns an empty string.
//
// Example:
//
//	input := `
//	    func main() {
//	        fmt.Println("Hello")
//	    }
//	`
//	result := Unindent(input)
//	// result is:
//	// func main() {
//	//     fmt.Println("Hello")
//	// }
func Unindent(s string) string {
	if s == "" {
		return ""
	}

	// Split into lines
	lines := strings.Split(s, "\n")

	// Remove leading blank lines
	for len(lines) > 0 && isBlankLine(lines[0]) {
		lines = lines[1:]
	}

	// Remove trailing blank lines
	for len(lines) > 0 && isBlankLine(lines[len(lines)-1]) {
		lines = lines[:len(lines)-1]
	}

	// Handle edge case: all lines were blank
	if len(lines) == 0 {
		return ""
	}

	// Find minimum indentation among non-empty lines
	minIndent := -1
	for _, line := range lines {
		if isBlankLine(line) {
			continue
		}
		indent := countLeadingWhitespace(line)
		if minIndent == -1 || indent < minIndent {
			minIndent = indent
		}
	}

	// If no non-empty lines found, return empty string
	if minIndent == -1 {
		return ""
	}

	// Remove common indentation and replace whitespace-only lines with blank lines
	for i, line := range lines {
		if isBlankLine(line) {
			lines[i] = ""
		} else if len(line) >= minIndent {
			lines[i] = line[minIndent:]
		}
	}

	// Join lines
	result := strings.Join(lines, "\n")

	// Ensure a trailing newline
	if !strings.HasSuffix(result, "\n") {
		result += "\n"
	}

	return result
}

// isBlankLine returns true if the line contains only whitespace or is empty
func isBlankLine(s string) bool {
	return strings.TrimSpace(s) == ""
}

// countLeadingWhitespace counts the number of leading whitespace characters
// (spaces and tabs, each counted as 1 character)
func countLeadingWhitespace(s string) int {
	count := 0
	for _, ch := range s {
		if ch == ' ' || ch == '\t' {
			count++
		} else {
			break
		}
	}
	return count
}
