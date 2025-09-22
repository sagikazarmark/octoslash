package builtin

import (
	"fmt"

	"github.com/google/go-github/v74/github"
	"github.com/spf13/cobra"
)

func NewCloseCommand(event github.IssueCommentEvent) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "close",
		Short: "Close a file",
		Long:  "Close a file",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("closing issue", event.GetIssue().GetNumber())
			return nil
		},
	}

	return cmd
}
