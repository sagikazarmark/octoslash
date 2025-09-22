package slash

import (
	"bufio"
	"bytes"
	"io"
	"strings"
)

// ScanString finds slash commands in the given string
func ScanString(text string) []string {
	return ScanReader(strings.NewReader(text))
}

// ScanBytes finds slash commands in the given byte slice
func ScanBytes(data []byte) []string {
	return ScanReader(bytes.NewReader(data))
}

// ScanReader reads from an io.Reader and extracts slash commands line by line
func ScanReader(r io.Reader) []string {
	var commands []string
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		line := scanner.Text()
		if cmd := extractCommandFromLine(line); cmd != "" {
			commands = append(commands, cmd)
		}
	}

	return commands
}

// extractCommandFromLine finds the first slash command in a single line
// Commands must be at the beginning of the line (after optional whitespace)
func extractCommandFromLine(line string) string {
	// Strip leading and trailing whitespace from the entire line
	line = strings.TrimSpace(line)

	// Check if line starts with a slash command
	if len(line) > 0 && line[0] == '/' {
		// Find the end of the command name
		end := 1
		for end < len(line) && !isWhitespace(line[end]) {
			end++
		}

		// Must have more than just '/' for a valid command
		if end > 1 {
			// Return the full command line without the leading slash
			return line[1:]
		}
	}

	return ""
}

// isWhitespace checks if a character is whitespace
func isWhitespace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\r'
}
