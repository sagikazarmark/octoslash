package builtin

import (
	"context"
	"log/slog"

	"github.com/google/go-github/v74/github"
	"github.com/spf13/cobra"
)

// Assign represents a command to assign an issue or pull request to a user.
type Assign struct {
	Repo     *github.Repository
	Issue    *github.Issue
	Assignee string
}

// AssignHandler handles the [Assign] command.
type AssignHandler struct {
	Client *github.Client
	Logger *slog.Logger
}

// Handle executes the [Assign] command.
func (h AssignHandler) Handle(ctx context.Context, cmd Assign) error {
	repo := cmd.Repo
	issue := cmd.Issue

	logger := h.Logger.With(slog.Int("number", issue.GetNumber()))

	logger.Info("assigning issue to user", slog.String("assignee", cmd.Assignee))

	_, _, err := h.Client.Issues.AddAssignees(
		ctx,
		repo.GetOwner().GetLogin(),
		repo.GetName(),
		issue.GetNumber(),
		[]string{cmd.Assignee},
	)
	if err != nil {
		return err
	}

	return nil
}

// NewAssignCommand creates a new Cobra command to assign an issue or pull request to a user.
//
// It integrates the [Assign] command into the default command dispatcher.
func NewAssignCommand(
	event github.IssueCommentEvent,
	client *github.Client,
	logger *slog.Logger,
) *cobra.Command {
	handler := AssignHandler{
		Client: client,
		Logger: logger,
	}

	return newAssignCommand(event, handler)
}

func newAssignCommand(
	event github.IssueCommentEvent,
	handler commandHandler[Assign],
) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "assign",
		Short: "Assign an issue or pull request to a user",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			command := Assign{
				Repo:     event.GetRepo(),
				Issue:    event.GetIssue(),
				Assignee: args[0],
			}

			return handler.Handle(cmd.Context(), command)
		},
	}

	return cmd
}

// SelfAssign represents a command to assign an issue or pull request to the current user.
type SelfAssign struct {
	Repo    *github.Repository
	Issue   *github.Issue
	Comment *github.IssueComment
}

// SelfAssignHandler handles the [SelfAssign] command.
type SelfAssignHandler struct {
	AssignHandler commandHandler[Assign]
}

// Handle executes the [SelfAssign] command.
func (h SelfAssignHandler) Handle(ctx context.Context, cmd SelfAssign) error {
	assign := Assign{
		Repo:     cmd.Repo,
		Issue:    cmd.Issue,
		Assignee: cmd.Comment.GetUser().GetLogin(),
	}

	return h.AssignHandler.Handle(ctx, assign)
}

// NewSelfAssignCommand creates a new Cobra command to assign an issue or pull request to the current user.
//
// It integrates the [SelfAssign] command into the default command dispatcher.
func NewSelfAssignCommand(
	event github.IssueCommentEvent,
	client *github.Client,
	logger *slog.Logger,
) *cobra.Command {
	handler := SelfAssignHandler{
		AssignHandler: AssignHandler{
			Client: client,
			Logger: logger,
		},
	}

	return newSelfAssignCommand(event, handler)
}

func newSelfAssignCommand(
	event github.IssueCommentEvent,
	handler commandHandler[SelfAssign],
) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "self-assign",
		Short: "Assign an issue or pull request to the current user",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			command := SelfAssign{
				Repo:    event.GetRepo(),
				Issue:   event.GetIssue(),
				Comment: event.GetComment(),
			}

			return handler.Handle(cmd.Context(), command)
		},
	}

	return cmd
}

// Unassign represents a command to unassign an issue or pull request from a user.
type Unassign struct {
	Repo     *github.Repository
	Issue    *github.Issue
	Assignee string
}

// UnassignHandler handles the [Unassign] command.
type UnassignHandler struct {
	Client *github.Client
	Logger *slog.Logger
}

// Handle executes the [Unassign] command.
func (h UnassignHandler) Handle(ctx context.Context, cmd Unassign) error {
	repo := cmd.Repo
	issue := cmd.Issue

	logger := h.Logger.With(slog.Int("number", issue.GetNumber()))

	logger.Info("unassigning issue from user", slog.String("assignee", cmd.Assignee))

	_, _, err := h.Client.Issues.RemoveAssignees(
		ctx,
		repo.GetOwner().GetLogin(),
		repo.GetName(),
		issue.GetNumber(),
		[]string{cmd.Assignee},
	)
	if err != nil {
		return err
	}

	return nil
}

// NewUnassignCommand creates a new Cobra command to unassign an issue or pull request from a user.
//
// It integrates the [Unassign] command into the default command dispatcher.
func NewUnassignCommand(
	event github.IssueCommentEvent,
	client *github.Client,
	logger *slog.Logger,
) *cobra.Command {
	handler := UnassignHandler{
		Client: client,
		Logger: logger,
	}

	return newUnassignCommand(event, handler)
}

func newUnassignCommand(
	event github.IssueCommentEvent,
	handler commandHandler[Unassign],
) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unassign",
		Short: "Unassign an issue or pull request from a user",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			command := Unassign{
				Repo:     event.GetRepo(),
				Issue:    event.GetIssue(),
				Assignee: args[0],
			}

			return handler.Handle(cmd.Context(), command)
		},
	}

	return cmd
}

// SelfUnassign represents a command to assign an issue or pull request to the current user.
type SelfUnassign struct {
	Repo    *github.Repository
	Issue   *github.Issue
	Comment *github.IssueComment
}

// SelfUnassignHandler handles the [SelfUnassign] command.
type SelfUnassignHandler struct {
	UnassignHandler commandHandler[Unassign]
}

// Handle executes the [SelfUnassign] command.
func (h SelfUnassignHandler) Handle(ctx context.Context, cmd SelfUnassign) error {
	unassign := Unassign{
		Repo:     cmd.Repo,
		Issue:    cmd.Issue,
		Assignee: cmd.Comment.GetUser().GetLogin(),
	}

	return h.UnassignHandler.Handle(ctx, unassign)
}

// NewSelfUnassignCommand creates a new Cobra command to unassign an issue or pull request from the current user.
//
// It integrates the [SelfUnassign] command into the default command dispatcher.
func NewSelfUnassignCommand(
	event github.IssueCommentEvent,
	client *github.Client,
	logger *slog.Logger,
) *cobra.Command {
	handler := SelfUnassignHandler{
		UnassignHandler: UnassignHandler{
			Client: client,
			Logger: logger,
		},
	}

	return newSelfUnassignCommand(event, handler)
}

func newSelfUnassignCommand(
	event github.IssueCommentEvent,
	handler commandHandler[SelfUnassign],
) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "self-unassign",
		Short: "Unassign an issue or pull request from the current user",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			command := SelfUnassign{
				Repo:    event.GetRepo(),
				Issue:   event.GetIssue(),
				Comment: event.GetComment(),
			}

			return handler.Handle(cmd.Context(), command)
		},
	}

	return cmd
}
