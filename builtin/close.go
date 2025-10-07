package builtin

import (
	"context"
	"log/slog"

	"github.com/google/go-github/v74/github"
	"github.com/spf13/cobra"
)

func NewCloseCommand(
	event github.IssueCommentEvent,
	client *github.Client,
	logger *slog.Logger,
) *cobra.Command {
	command := closeCommand{
		event:  event,
		client: client,
		logger: logger,
	}

	cmd := &cobra.Command{
		Use:   "close",
		Short: "Close an issue or pull request",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var reason string
			if len(args) > 0 {
				reason = args[0]
			}

			return command.run(cmd.Context(), reason)
		},
	}

	return cmd
}

type closeCommand struct {
	event  github.IssueCommentEvent
	client *github.Client
	logger *slog.Logger
}

func (c closeCommand) run(ctx context.Context, reason string) error {
	issue := c.event.GetIssue()
	repo := c.event.GetRepo()

	c.logger.Info("closing issue", slog.Int("number", issue.GetNumber()), slog.String("reason", ""))

	req := &github.IssueRequest{
		State: github.Ptr("closed"),
	}

	if reason != "" {
		req.StateReason = github.Ptr(reason)
	}

	_, _, err := c.client.Issues.Edit(
		ctx,
		repo.GetOwner().GetLogin(),
		repo.GetName(),
		issue.GetNumber(),
		req,
	)
	if err != nil {
		return err
	}

	return nil
}
