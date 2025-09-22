package slash_test

import (
	"strings"
	"testing"

	"github.com/sagikazarmark/octoslash/slash"
)

func TestParser_Parse(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected []string
		wantErr  bool
	}{
		{
			name:     "simple command",
			input:    "echo hello",
			expected: []string{"echo", "hello"},
		},
		{
			name:     "command with quoted string",
			input:    `echo "hello world"`,
			expected: []string{"echo", "hello world"},
		},
		{
			name:     "command with variables",
			input:    "echo $HOME $USER",
			expected: []string{"echo", "$HOME", "$USER"},
		},
		{
			name:     "command with flags",
			input:    "ls --foo bar --baz=qux",
			expected: []string{"ls", "--foo", "bar", "--baz=qux"},
		},
		{
			name:     "command with escaped characters",
			input:    `echo test\ file`,
			expected: []string{"echo", "test\\ file"},
		},
		{
			name:     "command with wildcards",
			input:    "ls * *.go",
			expected: []string{"ls", "*", "*.go"},
		},
		{
			name:     "command with command substitution",
			input:    `echo $(date)`,
			expected: []string{"echo", "$(date)"},
		},
		{
			name:     "complex command",
			input:    `echo "hello \"world\"" $HOME --foo bar --baz=$(pwd) test\ file`,
			expected: []string{"echo", "hello \"world\"", "$HOME", "--foo", "bar", "--baz=$(pwd)", "test\\ file"},
		},
		{
			name:    "empty input",
			input:   "",
			wantErr: true,
		},
		{
			name:    "whitespace only",
			input:   "   \t\n  ",
			wantErr: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			p := slash.NewParser()
			args, err := p.Parse(strings.NewReader(testCase.input))

			if testCase.wantErr {
				if err == nil {
					t.Errorf("expected error but got none")
				}

				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(args) != len(testCase.expected) {
				t.Fatalf("expected %d args, got %d: %v", len(testCase.expected), len(args), args)
			}

			for i, expected := range testCase.expected {
				if args[i] != expected {
					t.Errorf("arg %d: expected %q, got %q", i, expected, args[i])
				}
			}
		})
	}
}
