package builtin

import (
	"log/slog"

	"github.com/google/go-github/v74/github"
	"github.com/spf13/cobra"
)

func NewCloseCommand(client *github.Client, event github.IssueCommentEvent) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "close",
		Short: "Close a file",
		Long:  "Close a file",
		RunE: func(cmd *cobra.Command, args []string) error {
			slog.Info("closing issue", slog.Int("number", event.GetIssue().GetNumber()))

			_, _, err := client.Issues.Edit(cmd.Context(), event.GetRepo().GetOwner().GetLogin(), event.GetRepo().GetName(), event.GetIssue().GetNumber(), &github.IssueRequest{
				State: github.Ptr("closed"),
			})

			return err
		},
	}

	return cmd
}
