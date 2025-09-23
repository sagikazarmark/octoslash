package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/cedar-policy/cedar-go"
	"github.com/google/go-github/v74/github"
	githubfs "github.com/sagikazarmark/go-github-fs"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/sagikazarmark/octoslash/authz"
	"github.com/sagikazarmark/octoslash/builtin"
	"github.com/sagikazarmark/octoslash/slash"
)

func init() {
	cobra.EnableTraverseRunHooks = true
}

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	logger := slog.Default()

	flags := pflag.NewFlagSet("octoslash", pflag.ExitOnError)

	var eventName string
	flags.StringVar(&eventName, "event-name", os.Getenv("GITHUB_EVENT_NAME"), "")

	var eventPath string
	flags.StringVar(&eventPath, "event-path", os.Getenv("GITHUB_EVENT_PATH"), "")

	defaultConfigPath := filepath.Join(".github", "octoslash")
	if ws := os.Getenv("GITHUB_WORKSPACE"); ws != "" {
		defaultConfigPath = filepath.Join(ws, defaultConfigPath)
	}

	var configPath string
	flags.StringVar(&configPath, "config-path", defaultConfigPath, "")

	err := flags.Parse(os.Args[1:])
	if err != nil {
		logger.Error(err.Error())

		os.Exit(1)
	}

	logger = logger.With(slog.String("event_name", eventName), slog.String("event_path", eventPath))

	if eventName != "issue_comment" {
		logger.Error("unsupported event")

		os.Exit(1)
	}

	logger.Debug("loading event file")

	file, err := os.Open(eventPath)
	if err != nil {
		logger.Error(fmt.Sprintf("loading event: %s", err.Error()))

		os.Exit(1)
	}
	defer file.Close()

	logger.Debug("decoding event")

	var event github.IssueCommentEvent

	err = json.NewDecoder(file).Decode(&event)
	if err != nil {
		logger.Error(fmt.Sprintf("decoding event: %s", err.Error()))

		os.Exit(1)
	}

	fsys, err := openConfigPath(configPath)
	if err != nil {
		logger.Error(fmt.Sprintf("opening octoslash config: %s", err.Error()))

		os.Exit(1)
	}

	client := github.NewClient(nil)
	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		client = client.WithAuthToken(token)
	}

	// Try to use the GitHub filesystem
	if fsys == nil && os.Getenv("GITHUB_ACTIONS") == "true" {
		githubFS := githubfs.New(
			githubfs.WithClient(client),
			githubfs.WithRepository(event.GetRepo().GetOwner().GetLogin(), event.GetRepo().GetName()),
		)

		_, err := fs.Stat(githubFS, ".github/octoslash")
		if errors.Is(err, fs.ErrNotExist) {
			fsys = nil
		} else if err != nil {
			logger.Error(fmt.Sprintf("opening octoslash config from GitHub Actions: %s", err.Error()))

			os.Exit(1)
		} else {
			fsys = githubFS
		}
	}

	var policyLoaders authz.PolicyLoaders
	entityLoaders := authz.EntityLoaders{
		authz.EventEntityLoader{Event: event},
	}

	if fsys != nil {
		policyLoaders = append(policyLoaders, authz.FilePolicyLoader{
			Fsys: fsys,
		})
		entityLoaders = append(entityLoaders, authz.FileEntityLoader{
			Fsys: fsys,
		})
	} else {
		logger.Debug("no filesystem available, skipping policy and entity loading from filesystem")
	}

	policies, err := policyLoaders.LoadPolicies()
	if err != nil {
		logger.Error(fmt.Sprintf("loading policies: %s", err.Error()))

		os.Exit(1)
	}

	entities, err := entityLoaders.LoadEntities()
	if err != nil {
		logger.Error(fmt.Sprintf("loading entities: %s", err.Error()))

		os.Exit(1)
	}

	newCommand := func(ctx context.Context, args []string) *cobra.Command {
		rootCmd := &cobra.Command{
			Use:   "octoslash",
			Short: "A command-line tool for interacting with GitHub",
			Long:  "octoslash is a command-line tool for interacting with GitHub.",
			PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
				cmd.SilenceErrors = true
				cmd.SilenceUsage = true
				req := authz.NewRequest(cmd, event)

				ok, _ := cedar.Authorize(policies, entities, req)
				if !ok {
					return fmt.Errorf("principal %s is not authorized to perform %s on %s", req.Principal, req.Action, req.Resource)
				}

				return nil
			},
		}

		rootCmd.PersistentFlags().Bool("dry-run", false, "Do not perform any actions")

		rootCmd.AddCommand(builtin.NewCloseCommand(client, event))
		rootCmd.AddCommand(builtin.NewLabelCommand(client, event))
		rootCmd.AddCommand(builtin.NewRemoveLabelCommand(client, event))

		rootCmd.SetArgs(args)
		rootCmd.SetContext(ctx)

		// TODO: set output and error writers

		return rootCmd
	}

	// TODO: the scanner may return an error
	rawCommands := slash.ScanString(event.GetComment().GetBody())

	if len(rawCommands) == 0 {
		logger.Info("no commands to run")
	}

	parser := slash.NewParser()

	for _, rawCommand := range rawCommands {
		args, err := parser.Parse(strings.NewReader(rawCommand))
		if err != nil {
			logger.Error(fmt.Sprintf("parsing command: %s", err.Error()), slog.String("command", rawCommand))

			// TODO: make this behavior configurable
			continue
		}

		command := newCommand(context.Background(), args)

		if err := command.Execute(); err != nil {
			logger.Error(fmt.Sprintf("running command: %s", err.Error()), slog.String("command", rawCommand))

			os.Exit(1)
		}
	}
}

func openConfigPath(configPath string) (fs.FS, error) {
	root, err := os.OpenRoot(configPath)
	if err != nil && errors.Is(err, fs.ErrNotExist) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return root.FS(), nil
}
