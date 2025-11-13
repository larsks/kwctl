package helpers

import "testing"

func TestUnindent(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "single line with spaces",
			input:    "    hello",
			expected: "hello",
		},
		{
			name:     "single line with tabs",
			input:    "\t\thello",
			expected: "hello",
		},
		{
			name:     "single line without indentation",
			input:    "hello",
			expected: "hello",
		},
		{
			name:     "basic indented text",
			input:    "    line1\n    line2\n    line3",
			expected: "line1\nline2\nline3",
		},
		{
			name:     "preserve relative indentation",
			input:    "    line1\n        line2\n    line3",
			expected: "line1\n    line2\nline3",
		},
		{
			name:     "remove leading blank lines",
			input:    "\n\n    line1\n    line2",
			expected: "line1\nline2",
		},
		{
			name:     "remove trailing blank lines",
			input:    "    line1\n    line2\n\n\n",
			expected: "line1\nline2\n",
		},
		{
			name:     "remove leading and trailing blank lines",
			input:    "\n\n    line1\n    line2\n\n",
			expected: "line1\nline2\n",
		},
		{
			name:     "replace whitespace-only lines with blank",
			input:    "    line1\n        \n    line2",
			expected: "line1\n\nline2",
		},
		{
			name:     "mixed tabs and spaces",
			input:    "\t  line1\n\t  line2",
			expected: "line1\nline2",
		},
		{
			name:     "tabs and spaces with relative indentation",
			input:    "  \tline1\n  \t  line2\n  \tline3",
			expected: "line1\n  line2\nline3",
		},
		{
			name:     "preserve trailing newline",
			input:    "    line1\n    line2\n",
			expected: "line1\nline2\n",
		},
		{
			name:     "no trailing newline",
			input:    "    line1\n    line2",
			expected: "line1\nline2",
		},
		{
			name:     "all blank lines",
			input:    "\n\n\n",
			expected: "",
		},
		{
			name:     "all whitespace",
			input:    "    \n    \n    ",
			expected: "",
		},
		{
			name:     "only whitespace characters",
			input:    "    ",
			expected: "",
		},
		{
			name:     "no indentation needed",
			input:    "line1\nline2\nline3",
			expected: "line1\nline2\nline3",
		},
		{
			name:     "mixed indentation levels - find minimum",
			input:    "        line1\n    line2\n            line3",
			expected: "    line1\nline2\n        line3",
		},
		{
			name: "realistic code example",
			input: `
    func main() {
        fmt.Println("Hello")
        if true {
            fmt.Println("World")
        }
    }
`,
			expected: `func main() {
    fmt.Println("Hello")
    if true {
        fmt.Println("World")
    }
}
`,
		},
		{
			name:     "whitespace-only line in middle",
			input:    "    line1\n\t\t\n    line2",
			expected: "line1\n\nline2",
		},
		{
			name: "complex example with mixed whitespace",
			input: `

		line1
			line2
		line3

`,
			expected: "line1\n\tline2\nline3\n",
		},
		{
			name:     "single newline",
			input:    "\n",
			expected: "",
		},
		{
			name:     "text already at column 0",
			input:    "line1\n    line2\nline3",
			expected: "line1\n    line2\nline3",
		},
		{
			name:     "preserve empty lines between content",
			input:    "    line1\n\n    line2",
			expected: "line1\n\nline2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Unindent(tt.input)
			if result != tt.expected {
				t.Errorf("Unindent() = %q, expected %q", result, tt.expected)
			}
		})
	}
}
