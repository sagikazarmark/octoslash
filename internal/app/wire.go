package app

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"os"

	"github.com/google/go-github/v74/github"
	"github.com/wireinject/wire"

	githubfs "github.com/sagikazarmark/go-github-fs"
	"github.com/sagikazarmark/octoslash"
)

type Provider any

type Token string

type LocalFS = fs.FS

func InitializeEventHandler(
	provider Provider,
	token Token,
	repo *github.Repository,
	localFS LocalFS,
) (octoslash.EventHandler, error) {
	wire.Build(
		NewLogger,
		NewClient,
		NewFS,

		// Authorization
		DefaultAuthorizer,
		DefaultPolicyLoader,
		DefaultEntityLoader,

		NewCommandDispatcher,
		DefaultCommandDispatcher,
		NewCommandProvider,

		wire.Struct(new(octoslash.EventHandler), "*"),
	)

	return octoslash.EventHandler{}, nil
}

// TODO: inject output
func NewLogger(provider Provider) *slog.Logger {
	switch p := provider.(type) {
	case interface{ NewLogger() *slog.Logger }:
		return p.NewLogger()

	case interface{ NewLogger(io.Writer) *slog.Logger }:
		return p.NewLogger(os.Stderr)

	default:
		return slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug}))
	}
}

func NewClient(provider Provider, token Token) *github.Client {
	switch p := provider.(type) {
	case interface{ NewClient() *github.Client }:
		return p.NewClient()

	case interface{ NewClient(string) *github.Client }:
		return p.NewClient(string(token))

	default:
		client := github.NewClient(nil)
		if token != "" {
			client = client.WithAuthToken(string(token))
		}

		return client
	}
}

func NewFS(localFS LocalFS, client *github.Client, repo *github.Repository) LazyResult[fs.FS] {
	return func() (fs.FS, error) {
		if localFS != nil {
			return localFS, nil
		}

		githubFS := githubfs.New(
			githubfs.WithClient(client),
			githubfs.WithRepository(repo.GetOwner().GetLogin(), repo.GetName()),
		)

		const defaultConfigPath = ".github/octoslash"

		_, err := fs.Stat(githubFS, defaultConfigPath)
		if errors.Is(err, fs.ErrNotExist) {
			return nil, nil
		} else if err != nil {
			return nil, fmt.Errorf("opening octoslash config from GitHub Actions: %w", err)
		}

		return fs.Sub(githubFS, defaultConfigPath)
	}
}
