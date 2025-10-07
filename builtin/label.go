package builtin

import (
	"context"
	"log/slog"

	"github.com/google/go-github/v74/github"
	"github.com/spf13/cobra"
)

func NewLabelCommand(event github.IssueCommentEvent, client *github.Client, logger *slog.Logger) *cobra.Command {
	command := labelCommand{
		event:  event,
		client: client,
		logger: logger,
	}

	cmd := &cobra.Command{
		Use:   "label",
		Short: "Label an issue or pull request",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return command.run(cmd.Context(), args[0])
		},
	}

	return cmd
}

type labelCommand struct {
	event  github.IssueCommentEvent
	client *github.Client
	logger *slog.Logger
}

func (c labelCommand) run(ctx context.Context, label string) error {
	issue := c.event.GetIssue()
	repo := c.event.GetRepo()

	logger := c.logger.With(slog.Int("number", issue.GetNumber()))

	logger.Info("adding label to issue", slog.String("label", label))

	_, _, err := c.client.Issues.AddLabelsToIssue(
		ctx,
		repo.GetOwner().GetLogin(),
		repo.GetName(),
		issue.GetNumber(),
		[]string{label},
	)
	if err != nil {
		return err
	}

	return nil
}

func NewRemoveLabelCommand(event github.IssueCommentEvent, client *github.Client, logger *slog.Logger) *cobra.Command {
	command := removeLabelCommand{
		event:  event,
		client: client,
		logger: logger,
	}

	cmd := &cobra.Command{
		Use:   "remove-label",
		Short: "Remove a label from an issue or pull request",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return command.run(cmd.Context(), args[0])
		},
	}

	return cmd
}

type removeLabelCommand struct {
	event  github.IssueCommentEvent
	client *github.Client
	logger *slog.Logger
}

func (c removeLabelCommand) run(ctx context.Context, label string) error {
	issue := c.event.GetIssue()
	repo := c.event.GetRepo()

	logger := c.logger.With(
		slog.Int("number", issue.GetNumber()),
	)

	logger.Info("removing label from issue", slog.String("label", label))

	_, err := c.client.Issues.RemoveLabelForIssue(
		ctx,
		repo.GetOwner().GetLogin(),
		repo.GetName(),
		issue.GetNumber(),
		label,
	)
	if err != nil {
		return err
	}

	return nil
}
