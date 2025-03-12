package stdin

import (
	"bytes"
	"testing"
)

// H// Test implementation that bypasses os.Stdin completely for more reliable testing
func TestRead(t *testing.T) {

	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Empty input",
			input:    "",
			expected: "",
		},
		{
			name:     "Single line input",
			input:    "Hello, world!",
			expected: "Hello, world!",
		},
		{
			name:     "Multi-line input",
			input:    "Line 1\nLine 2\nLine 3",
			expected: "Line 1\nLine 2\nLine 3",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a buffer with our test input
			buf := bytes.NewBufferString(tc.input)

			// Test with our buffer directly, avoiding os.Stdin
			result, err := Read(buf)

			// Check for errors
			if err != nil {
				t.Fatalf("readFromReader returned an error: %v", err)
			}

			// Check the result
			if result != tc.expected {
				t.Errorf("Expected %q, got %q", tc.expected, result)
			}
		})
	}
}
