package builtin

import (
	"fmt"
	"log/slog"

	"github.com/google/go-github/v74/github"
	"github.com/spf13/cobra"
)

func NewLabelCommand(client *github.Client, event github.IssueCommentEvent) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "label",
		Short: "Label an issue",
		Long:  "Label an issue",
		RunE: func(cmd *cobra.Command, args []string) error {
			slog.Info("adding label", slog.Int("number", event.GetIssue().GetNumber()))

			if len(args) == 0 {
				return fmt.Errorf("no labels provided")
			}

			_, _, err := client.Issues.AddLabelsToIssue(cmd.Context(), event.GetRepo().GetOwner().GetLogin(), event.GetRepo().GetName(), event.GetIssue().GetNumber(), []string{args[0]})
			if err != nil {
				return err
			}

			return nil
		},
	}

	return cmd
}

func NewRemoveLabelCommand(client *github.Client, event github.IssueCommentEvent) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove-label",
		Short: "Remove label from an issue",
		Long:  "Remove label from an issue",
		RunE: func(cmd *cobra.Command, args []string) error {
			slog.Info("removing label", slog.Int("number", event.GetIssue().GetNumber()))

			if len(args) == 0 {
				return fmt.Errorf("no labels provided")
			}

			_, err := client.Issues.RemoveLabelForIssue(cmd.Context(), event.GetRepo().GetOwner().GetLogin(), event.GetRepo().GetName(), event.GetIssue().GetNumber(), args[0])
			if err != nil {
				return err
			}

			return nil
		},
	}

	return cmd
}
