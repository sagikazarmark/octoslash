package command

import (
	"bytes"
	"context"
	"errors"
	"io"

	"github.com/google/go-github/v74/github"
	"github.com/spf13/cobra"
)

type CobraDispatcher struct {
	Authorizer      Authorizer
	CommandProvider CommandProvider
}

type Authorizer interface {
	Authorize(ctx context.Context, event github.IssueCommentEvent, action string) error
}

type CommandProvider interface {
	NewCommand(event github.IssueCommentEvent) *cobra.Command
}

func (d CobraDispatcher) Dispatch(
	ctx context.Context,
	event github.IssueCommentEvent,
	args []string,
) error {
	cmd := d.newCommand(event, args)

	return cmd.ExecuteContext(ctx)
}

func (d CobraDispatcher) newCommand(event github.IssueCommentEvent, args []string) *cobra.Command {
	cmd := d.CommandProvider.NewCommand(event)

	prevPersistentPreRunE := cmd.PersistentPreRunE

	authorizer := d.Authorizer

	cmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		if authorizer == nil {
			return errors.New("no authorizer configured, denying request")
		}

		action := cmd.Name()

		rootCmd := cmd.Root()
		cmd.VisitParents(func(cmd *cobra.Command) {
			if cmd == rootCmd {
				return
			}

			action = cmd.Name() + ":" + action
		})

		err := authorizer.Authorize(cmd.Context(), event, action)
		if err != nil {
			return err
		}

		if prevPersistentPreRunE != nil {
			return prevPersistentPreRunE(cmd, args)
		}

		return nil
	}

	// TODO: set output and error writers

	cmd.SilenceErrors = true
	cmd.SilenceUsage = true

	// TODO: should these be discarded?
	cmd.SetOut(io.Discard)
	cmd.SetErr(io.Discard)
	cmd.SetIn(&bytes.Buffer{})

	cmd.SetArgs(args)

	return cmd
}
