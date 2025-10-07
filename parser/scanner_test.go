package parser

import (
	"reflect"
	"strings"
	"testing"
)

// stringSlicesEqual compares two string slices, treating nil and empty slices as equal
func stringSlicesEqual(a, b []string) bool {
	if len(a) == 0 && len(b) == 0 {
		return true
	}
	return reflect.DeepEqual(a, b)
}

func TestScanString(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		expected []string
	}{
		{
			name:     "single command at start of line",
			text:     "/help",
			expected: []string{"help"},
		},
		{
			name:     "command with space before",
			text:     " /status",
			expected: []string{"status"},
		},
		{
			name:     "only first command on line is detected",
			text:     "/help /status",
			expected: []string{"help /status"},
		},
		{
			name:     "commands on multiple lines",
			text:     "/help\n/status\n/quit",
			expected: []string{"help", "status", "quit"},
		},
		{
			name:     "command with text after",
			text:     "/help me with this",
			expected: []string{"help me with this"},
		},
		{
			name:     "text before the command (should not match)",
			text:     "Please use /help for assistance",
			expected: []string{},
		},
		{
			name:     "slash in the middle of a line",
			text:     "The file is in home/user/documents",
			expected: []string{},
		},
		{
			name:     "no slash commands",
			text:     "hello world\nthis is text",
			expected: []string{},
		},
		{
			name:     "slash in middle of word (not a command)",
			text:     "http://example.com",
			expected: []string{},
		},
		{
			name:     "just slash character",
			text:     "/",
			expected: []string{},
		},
		{
			name:     "command with tab before",
			text:     "\t/command",
			expected: []string{"command"},
		},
		{
			name:     "empty text",
			text:     "",
			expected: []string{},
		},
		{
			name:     "mixed content with commands (only line-start commands match)",
			text:     "Hello /world\nThis is a test /check this out\n/final",
			expected: []string{"final"},
		},
		{
			name:     "command at end of sentence (should not match)",
			text:     "Type /quit to exit.",
			expected: []string{},
		},
		{
			name:     "multiple lines with text before commands (should not match)",
			text:     "Use /help for help\nTry /status for status\nFinally /quit to exit",
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ScanString(tt.text)
			if !stringSlicesEqual(result, tt.expected) {
				t.Errorf("ScanString() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestScanBytes(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		expected []string
	}{
		{
			name:     "single command",
			data:     []byte("/help"),
			expected: []string{"help"},
		},
		{
			name:     "multiple commands",
			data:     []byte("/help\n/status\n/quit"),
			expected: []string{"help", "status", "quit"},
		},
		{
			name:     "mixed content",
			data:     []byte("Hello world\n/help\nRegular text\n  /status"),
			expected: []string{"help", "status"},
		},
		{
			name:     "no commands",
			data:     []byte("Hello world\nThis is text"),
			expected: []string{},
		},
		{
			name:     "empty data",
			data:     []byte{},
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ScanBytes(tt.data)
			if !stringSlicesEqual(result, tt.expected) {
				t.Errorf("ScanBytes() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestScanReader(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "single command",
			input:    "/help",
			expected: []string{"help"},
		},
		{
			name:     "multiple lines with commands",
			input:    "/help\n/status\n/quit",
			expected: []string{"help", "status", "quit"},
		},
		{
			name:     "mixed content",
			input:    "Hello world\n/help\nRegular text\n  /status",
			expected: []string{"help", "status"},
		},
		{
			name:     "no commands",
			input:    "Hello world\nThis is text",
			expected: []string{},
		},
		{
			name:     "empty input",
			input:    "",
			expected: []string{},
		},
		{
			name:     "commands with text before (should not match)",
			input:    "Please use /help\nTry /status for info",
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.input)
			result := ScanReader(reader)
			if !stringSlicesEqual(result, tt.expected) {
				t.Errorf("ScanReader() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestExtractCommandFromLine(t *testing.T) {
	tests := []struct {
		name     string
		line     string
		expected string
	}{
		{
			name:     "single command",
			line:     "/test",
			expected: "test",
		},
		{
			name:     "no commands",
			line:     "regular text",
			expected: "",
		},
		{
			name:     "first command only",
			line:     "/first /second",
			expected: "first /second",
		},
		{
			name:     "command with text before (should not match)",
			line:     "Please use /help",
			expected: "",
		},
		{
			name:     "slash in middle of word",
			line:     "path/to/file",
			expected: "",
		},
		{
			name:     "command with whitespace before",
			line:     "  /status",
			expected: "status",
		},
		{
			name:     "just slash",
			line:     "/",
			expected: "",
		},
		{
			name:     "command with arguments",
			line:     "/help me with this",
			expected: "help me with this",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractCommandFromLine(tt.line)
			if result != tt.expected {
				t.Errorf("extractCommandFromLine() = %v, expected %v", result, tt.expected)
			}
		})
	}
}
