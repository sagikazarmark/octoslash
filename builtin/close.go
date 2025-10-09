package builtin

import (
	"context"
	"log/slog"

	"github.com/google/go-github/v74/github"
	"github.com/spf13/cobra"
)

// Close represents a command to close an issue or pull request.
type Close struct {
	Repo   *github.Repository
	Issue  *github.Issue
	Reason string
}

// CloseHandler handles the [Close] command.
type CloseHandler struct {
	Client *github.Client
	Logger *slog.Logger
}

// Handle executes the [Close] command.
func (h CloseHandler) Handle(ctx context.Context, cmd Close) error {
	issue := cmd.Issue
	repo := cmd.Repo

	h.Logger.Info(
		"closing issue",
		slog.Int("number", issue.GetNumber()),
		slog.String("reason", cmd.Reason),
	)

	req := &github.IssueRequest{
		State: github.Ptr("closed"),
	}

	if cmd.Reason != "" {
		req.StateReason = github.Ptr(cmd.Reason)
	}

	_, _, err := h.Client.Issues.Edit(
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

// NewCloseCommand creates a new Cobra command to close an issue or pull request.
//
// It integrates the [Close] command into the default command dispatcher.
func NewCloseCommand(
	event github.IssueCommentEvent,
	client *github.Client,
	logger *slog.Logger,
) *cobra.Command {
	handler := CloseHandler{
		Client: client,
		Logger: logger,
	}

	return newCloseCommand(event, handler)
}

func newCloseCommand(
	event github.IssueCommentEvent,
	handler commandHandler[Close],
) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "close",
		Short: "Close an issue or pull request",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var reason string
			if len(args) > 0 {
				reason = args[0]
			}

			command := Close{
				Repo:   event.GetRepo(),
				Issue:  event.GetIssue(),
				Reason: reason,
			}

			return handler.Handle(cmd.Context(), command)
		},
	}

	return cmd
}
