package builtin

import (
	"context"
	"log/slog"

	"github.com/google/go-github/v74/github"
	"github.com/spf13/cobra"
)

// AddLabel represents a command to add a label to an issue or pull request.
type AddLabel struct {
	Repo  *github.Repository
	Issue *github.Issue
	Label string
}

// AddLabelHandler handles the [AddLabel] command.
type AddLabelHandler struct {
	Client *github.Client
	Logger *slog.Logger
}

// Handle executes the [AddLabel] command.
func (h AddLabelHandler) Handle(ctx context.Context, cmd AddLabel) error {
	repo := cmd.Repo
	issue := cmd.Issue

	logger := h.Logger.With(slog.Int("number", issue.GetNumber()))

	logger.Info("adding label to issue", slog.String("label", cmd.Label))

	_, _, err := h.Client.Issues.AddLabelsToIssue(
		ctx,
		repo.GetOwner().GetLogin(),
		repo.GetName(),
		issue.GetNumber(),
		[]string{cmd.Label},
	)
	if err != nil {
		return err
	}

	return nil
}

// NewAddLabelCommand creates a new Cobra command to add a label to an issue or pull request.
//
// It integrates the [AddLabel] command into the default command dispatcher.
func NewAddLabelCommand(
	event github.IssueCommentEvent,
	client *github.Client,
	logger *slog.Logger,
) *cobra.Command {
	handler := AddLabelHandler{
		Client: client,
		Logger: logger,
	}

	return newAddLabelCommand(event, handler)
}

func newAddLabelCommand(
	event github.IssueCommentEvent,
	handler commandHandler[AddLabel],
) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "add-label",
		Aliases: []string{"label"},
		Short:   "Label an issue or pull request",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			command := AddLabel{
				Repo:  event.GetRepo(),
				Issue: event.GetIssue(),
				Label: args[0],
			}

			return handler.Handle(cmd.Context(), command)
		},
	}

	return cmd
}

// RemoveLabel represents a command to remove a label from an issue or pull request.
type RemoveLabel struct {
	Repo  *github.Repository
	Issue *github.Issue
	Label string
}

// RemoveLabelHandler handles the [RemoveLabel] command.
type RemoveLabelHandler struct {
	Client *github.Client
	Logger *slog.Logger
}

// Handle executes the [RemoveLabel] command.
func (h RemoveLabelHandler) Handle(ctx context.Context, cmd RemoveLabel) error {
	issue := cmd.Issue
	repo := cmd.Repo

	logger := h.Logger.With(slog.Int("number", issue.GetNumber()))

	logger.Info("removing label from issue", slog.String("label", cmd.Label))

	_, err := h.Client.Issues.RemoveLabelForIssue(
		ctx,
		repo.GetOwner().GetLogin(),
		repo.GetName(),
		issue.GetNumber(),
		cmd.Label,
	)
	if err != nil {
		return err
	}

	return nil
}

// NewRemoveLabelCommand creates a new Cobra command to remove a label from an issue or pull request.
//
// It integrates the [RemoveLabel] command into the default command dispatcher.
func NewRemoveLabelCommand(
	event github.IssueCommentEvent,
	client *github.Client,
	logger *slog.Logger,
) *cobra.Command {
	handler := RemoveLabelHandler{
		Client: client,
		Logger: logger,
	}

	return newRemoveLabelCommand(event, handler)
}

func newRemoveLabelCommand(
	event github.IssueCommentEvent,
	handler commandHandler[RemoveLabel],
) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove-label",
		Short: "Remove a label from an issue or pull request",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			command := RemoveLabel{
				Repo:  event.GetRepo(),
				Issue: event.GetIssue(),
				Label: args[0],
			}

			return handler.Handle(cmd.Context(), command)
		},
	}

	return cmd
}
