package builtin

import (
	"log/slog"

	"github.com/google/go-github/v74/github"
	"github.com/spf13/cobra"
)

type Provider struct{}

func (p Provider) NewCommandProvider(client *github.Client, logger *slog.Logger) CommandProvider {
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
		NewLabelCommand(event, p.Client, p.Logger),
		NewRemoveLabelCommand(event, p.Client, p.Logger),
	)

	return rootCmd
}
