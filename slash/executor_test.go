package slash_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/sagikazarmark/octoslash/slash"
	"github.com/spf13/cobra"
)

func TestExecutor_Execute(t *testing.T) {
	var executed []string

	newCommand := func() *cobra.Command {
		rootCmd := &cobra.Command{
			Use: "test",
			Run: func(cmd *cobra.Command, args []string) {
				executed = append(executed, fmt.Sprintf("%s %v", cmd.Use, args))
			},
		}

		deployCmd := &cobra.Command{
			Use: "deploy",
			Run: func(cmd *cobra.Command, args []string) {
				executed = append(executed, fmt.Sprintf("deploy %v", args))
			},
		}

		statusCmd := &cobra.Command{
			Use: "status",
			Run: func(cmd *cobra.Command, args []string) {
				executed = append(executed, fmt.Sprintf("status %v", args))
			},
		}

		rootCmd.AddCommand(deployCmd)
		rootCmd.AddCommand(statusCmd)
		return rootCmd
	}

	testCases := []struct {
		name     string
		input    string
		expected []string
		wantErr  bool
	}{
		{
			name:     "single command",
			input:    "/deploy",
			expected: []string{"deploy []"},
		},
		{
			name:     "multiple commands",
			input:    "/deploy\n/status",
			expected: []string{"deploy []", "status []"},
		},
		{
			name:     "empty input",
			input:    "",
			expected: nil,
		},
		{
			name:     "command with arguments",
			input:    "/deploy app",
			expected: []string{"deploy [app]"},
		},
		{
			name:     "command with quoted argument",
			input:    "/deploy \"my app\"",
			expected: []string{"deploy [my app]"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			executed = nil
			executor := slash.Executor{
				NewCommand: newCommand,
			}

			err := executor.Execute(context.Background(), tc.input)

			if tc.wantErr && err == nil {
				t.Error("expected error but got none")
			}
			if !tc.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if len(executed) != len(tc.expected) {
				t.Errorf("expected %d commands executed, got %d", len(tc.expected), len(executed))
			}

			for i, exp := range tc.expected {
				if i >= len(executed) || executed[i] != exp {
					t.Errorf("expected command %d to be %q, got %q", i, exp, executed[i])
				}
			}
		})
	}
}
