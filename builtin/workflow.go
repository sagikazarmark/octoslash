package builtin

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/google/go-github/v74/github"
	"github.com/spf13/cobra"
)

// WorkflowRun represents a command to run a workflow that accepts a workflow_dispatch trigger on a pull request.
type WorkflowRun struct {
	Repo             *github.Repository
	Issue            *github.Issue
	WorkflowFileName string
	Inputs           map[string]any
}

// WorkflowRunHandler handles the [WorkflowRun] command.
type WorkflowRunHandler struct {
	Client *github.Client
	Logger *slog.Logger
}

// Handle executes the [Assign] command.
func (h WorkflowRunHandler) Handle(ctx context.Context, cmd WorkflowRun) error {
	repo := cmd.Repo
	issue := cmd.Issue

	logger := h.Logger.With(slog.Int("number", issue.GetNumber()))

	if !issue.IsPullRequest() {
		return errors.New("cannot run workflow for issues")
	}

	logger.Info("fetching pull request details")

	pr, _, err := h.Client.PullRequests.Get(
		ctx,
		repo.GetOwner().GetLogin(),
		repo.GetName(),
		issue.GetNumber(),
	)
	if err != nil {
		return err
	}

	logger.Info("running workflow", slog.String("workflow", cmd.WorkflowFileName))

	_, err = h.Client.Actions.CreateWorkflowDispatchEventByFileName(
		ctx,
		repo.GetOwner().GetLogin(),
		repo.GetName(),
		cmd.WorkflowFileName,
		github.CreateWorkflowDispatchEventRequest{
			Ref:    pr.GetHead().GetRef(),
			Inputs: cmd.Inputs,
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// NewWorkflowRunCommand creates a new Cobra command to run a workflow on a pull request.
//
// It integrates the [WorkflowRun] command into the default command dispatcher.
func NewWorkflowRunCommand(
	event github.IssueCommentEvent,
	client *github.Client,
	logger *slog.Logger,
) *cobra.Command {
	handler := WorkflowRunHandler{
		Client: client,
		Logger: logger,
	}

	return newWorkflowRunCommand(event, handler)
}

func newWorkflowRunCommand(
	event github.IssueCommentEvent,
	handler commandHandler[WorkflowRun],
) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "workflow-run",
		Short: "Run a workflow on a pull request",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			inputs := make(map[string]any, len(args)-1)

			for _, arg := range args[1:] {
				keyVal := strings.SplitN(arg, "=", 2)
				if len(keyVal) != 2 {
					return fmt.Errorf("invalid input: %q", arg)
				}

				inputs[keyVal[0]] = keyVal[1]
			}

			command := WorkflowRun{
				Repo:             event.GetRepo(),
				Issue:            event.GetIssue(),
				WorkflowFileName: args[0],
				Inputs:           inputs,
			}

			return handler.Handle(cmd.Context(), command)
		},
	}

	return cmd
}
