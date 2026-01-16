package hash

import (
	"testing"
)

func TestSha256(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "numeric string",
			input:    "1234",
			expected: "03ac674216f3e15c761ee1a5e255f067953623c8b388b4459e13f978d7c846f4",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		},
		{
			name:     "hello world",
			input:    "hello world",
			expected: "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9",
		},
		{
			name:     "unicode",
			input:    "日本語",
			expected: "77710aedc74ecfa33685e33a6c7df5cc83004da1bdcef7fb280f5c2b2e97e0a5",
		},
		{
			name:     "with spaces",
			input:    "  spaces  ",
			expected: "3294fa4c529bfcd6d2b5683a79ea6ca944a92bc2dd8b2601cb2427a6abd0ba04",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Sha256(tt.input)
			if result != tt.expected {
				t.Errorf("Sha256(%q) = %q, want %q", tt.input, result, tt.expected)
			}
			// Verify output format: 64 lowercase hex characters
			if len(result) != 64 {
				t.Errorf("Sha256(%q) length = %d, want 64", tt.input, len(result))
			}
		})
	}
}
