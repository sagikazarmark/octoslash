package builtin

import (
	"github.com/google/go-github/v74/github"
	"github.com/spf13/cobra"
)

func NewLabelCommand(event github.IssueCommentEvent) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "label",
		Short: "Label a file",
		Long:  "Label a file",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	return cmd
}
