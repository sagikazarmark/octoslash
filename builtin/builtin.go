package builtin

import (
	"context"
	"log/slog"

	"github.com/google/go-github/v74/github"
	"github.com/spf13/cobra"

	"github.com/sagikazarmark/octoslash/command"
)

type Provider struct{}

func (p Provider) NewCommandProvider(
	client *github.Client,
	logger *slog.Logger,
) command.CommandProvider {
	return CommandProvider{
		Client: client,
		Logger: logger,
	}
}

type CommandProvider struct {
	Client *github.Client
	Logger *slog.Logger
}

func (p CommandProvider) NewCommand(event github.IssueCommentEvent) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "octoslash",
		Short: "Slash commands for GitHub issues and pull requests",
	}

	rootCmd.AddCommand(
		NewCloseCommand(event, p.Client, p.Logger),

		NewAddLabelCommand(event, p.Client, p.Logger),
		NewRemoveLabelCommand(event, p.Client, p.Logger),

		NewAssignCommand(event, p.Client, p.Logger),
		NewSelfAssignCommand(event, p.Client, p.Logger),
		NewUnassignCommand(event, p.Client, p.Logger),
		NewSelfUnassignCommand(event, p.Client, p.Logger),

		NewWorkflowRunCommand(event, p.Client, p.Logger),
	)

	return rootCmd
}

type commandHandler[T any] interface {
	Handle(ctx context.Context, command T) error
}
