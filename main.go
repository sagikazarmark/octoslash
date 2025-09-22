package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
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
	logger := slog.Default()

	flags := pflag.NewFlagSet("octoslash", pflag.ExitOnError)

	var eventName string
	flags.StringVar(&eventName, "event-name", os.Getenv("GITHUB_EVENT_NAME"), "")

	var eventPath string
	flags.StringVar(&eventPath, "event-path", os.Getenv("GITHUB_EVENT_PATH"), "")

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

	isGitHubActions := os.Getenv("GITHUB_ACTIONS") == "true"

	var fsys fs.FS

	if isGitHubActions {
		logger.Debug("running in GitHub Actions")

		fsys = githubfs.New()

		// Try to load GITHUB_WORKSPACE/.github/octoslash
		// If it doesn't exist, try to use githubfs
	} else {
		root, err := os.OpenRoot(".github/octoslash")
		if err != nil {
			logger.Error(fmt.Sprintf("opening octoslash config: %s", err.Error()))

			os.Exit(1)
		}

		fsys = root.FS()
	}

	loader := authz.FilePolicyLoader{
		Fsys: fsys,
	}

	policies, err := loader.LoadPolicies()
	if err != nil {
		panic(err)
	}

	e, err := entities(event, fsys)
	if err != nil {
		panic(err)
	}

	newCommand := func(ctx context.Context, args []string) *cobra.Command {
		rootCmd := &cobra.Command{
			Use:   "octoslash",
			Short: "A command-line tool for interacting with GitHub",
			Long:  "octoslash is a command-line tool for interacting with GitHub.",
			PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
				// cmd.SilenceErrors = true
				cmd.SilenceUsage = true
				req := authz.NewRequest(cmd, event)

				ok, _ := cedar.Authorize(policies, e, req)
				if !ok {
					return fmt.Errorf("principal %s is not authorized to perform %s on %s", req.Principal, req.Action, req.Resource)
				}

				return nil
			},
		}

		rootCmd.AddCommand(builtin.NewCloseCommand(event))
		rootCmd.AddCommand(builtin.NewLabelCommand(event))

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

func entities(event github.IssueCommentEvent, fsys fs.FS) (cedar.EntityGetter, error) {
	entities, err := loadPrincipals(fsys)
	if err != nil {
		return nil, err
	}

	owner := authz.NewOwner(event.GetRepo().GetOwner())
	repo := authz.NewRepository(event.GetRepo())
	issue := authz.NewIssueOrPullRequest(event.GetIssue(), event.GetRepo())

	entities[owner.UID] = owner
	entities[repo.UID] = repo
	entities[issue.UID] = issue

	return entities, nil
}

func loadPrincipals(fsys fs.FS) (cedar.EntityMap, error) {
	var entities cedar.EntityMap

	file, err := fsys.Open("principals.json")
	if errors.Is(err, fs.ErrNotExist) {
		return entities, nil
	}
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&entities); err != nil {
		return nil, err
	}

	return entities, nil
}
